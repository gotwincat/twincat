// Copyright 2021 gotwincat authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

// Package twincat implements a Twincat v3 client
// for the Beckhoff AMS protocol.
package twincat

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gotwincat/twincat/ams"
)

var ErrTimeout = errors.New("timeout")

// Client implements a Twincat3 TCP client.
type Client struct {
	Addr        string
	ReadTimeout time.Duration

	conn         net.Conn
	nextInvokeID uint32 // atomic

	mu      sync.Mutex
	handler map[uint32]chan ams.Response

	adsState    atomic.Value // uint16
	deviceState atomic.Value // uint16
}

func (c *Client) ADSState() uint16 {
	return c.adsState.Load().(uint16)
}

func (c *Client) SetADSState(s uint16) {
	c.adsState.Store(s)
}

func (c *Client) DeviceState() uint16 {
	return c.deviceState.Load().(uint16)
}

func (c *Client) SetDeviceState(s uint16) {
	c.deviceState.Store(s)
}

// Dial connects to a Twincat server.
func (c *Client) Dial(ctx context.Context) error {
	atomic.AddUint32(&c.nextInvokeID, 1)

	c.SetADSState(ams.ADSStateStart)
	c.SetDeviceState(ams.ADSStateStart)

	d := &net.Dialer{}
	conn, err := d.DialContext(ctx, "tcp", c.Addr)
	if err != nil {
		return err
	}
	c.conn = conn
	go c.receive(ctx)
	return nil
}

func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

type responseDecoder interface {
	ams.Response
	ams.Decoder
}

func (c *Client) receive(ctx context.Context) error {
	c.SetADSState(ams.ADSStateRun)
	c.SetDeviceState(ams.ADSStateRun)
	defer c.SetADSState(ams.ADSStateStop)
	defer c.SetDeviceState(ams.ADSStateStop)

	// We assume that a packet fits into a single packet.
	// This is probably wrong but I haven't found anything on length
	// would probably have to read the header first, alloc and then read
	// the rest to fix this probably. This works for now.
	const packetSize = 1500

	for {
		// read the next packet
		data := make([]byte, packetSize)
		n, err := c.conn.Read(data)
		if err != nil {
			return err
		}

		// truncate the buffer to the correct length
		data = data[:n]
		// log.Printf("read %d bytes", n)

		// decode just the header
		var hdr ams.Header
		if err := hdr.Decode(ams.NewBuffer(data)); err != nil {
			return err
		}

		// figure out the packet type
		var pkt packet
		switch {
		case ams.IsReadResponse(hdr.AMSHeader):
			pkt = &ams.ReadResponse{}
		case ams.IsWriteResponse(hdr.AMSHeader):
			pkt = &ams.WriteResponse{}
		case ams.IsReadWriteResponse(hdr.AMSHeader):
			pkt = &ams.ReadWriteResponse{}
		case ams.IsReadStateRequest(hdr.AMSHeader):
			pkt = &ams.ReadStateRequest{}
		default:
			log.Printf("client: unknown packet: %#v", hdr)
			continue
		}

		// decode the full packet with the header
		if err := pkt.Decode(ams.NewBuffer(data)); err != nil {
			log.Printf("client: failed to decode: %s", err)
			return err
		}

		switch req := pkt.(type) {
		// handle incoming requests
		case *ams.ReadStateRequest:
			if err := c.handleReadStateRequest(ctx, req); err != nil {
				return err
			}

		// forward responses to handlers
		default:
			// find the handler channel for packet
			invokeID := hdr.AMSHeader.InvokeID
			c.mu.Lock()
			if c.handler == nil {
				c.handler = make(map[uint32]chan ams.Response)
			}
			h := c.handler[invokeID]
			delete(c.handler, invokeID)
			c.mu.Unlock()

			// if there is no handler then drop the packet
			if h == nil {
				log.Printf("client: no handler for %d", invokeID)
				continue
			}

			// otherwise send the response to the handler.
			// here we assume that h is buffered and can hold
			// one response. So this call should never block.
			select {
			case <-ctx.Done():
			case h <- pkt:
				close(h)
			}
		}
	}
}

func (c *Client) handleReadStateRequest(ctx context.Context, req *ams.ReadStateRequest) error {
	hdr := req.Header()
	resp := ams.NewReadStateResponse(hdr.Sender, hdr.Target, ams.NoError, c.ADSState(), c.DeviceState())
	return c.sendResponse(ctx, req, resp)
}

type packet interface {
	Header() *ams.AMSHeader
	Decode(b *ams.Buffer) error
	Encode(b *ams.Buffer) error
}

func (c *Client) sendResponse(ctx context.Context, req ams.Request, pkt packet) error {
	// set the invoke id from the request
	pkt.Header().InvokeID = req.Header().InvokeID

	// encode the response
	var b ams.Buffer
	if err := pkt.Encode(&b); err != nil {
		return err
	}

	// send the response
	_, err := c.conn.Write(b.Bytes())
	return err
}

// send sends a request to the server and sets up a handler channel
// for the callback.
func (c *Client) send(ctx context.Context, pkt packet, cb func(ams.Response) error) error {
	// set a unique invoke id for the request
	pkt.Header().InvokeID = atomic.AddUint32(&c.nextInvokeID, 1)

	// encode the request
	var b ams.Buffer
	if err := pkt.Encode(&b); err != nil {
		return err
	}

	// create a handler channel for the response
	// make sure that the channel is buffered
	// so that we don't need a separate go routine for
	// sending the resposne.
	h := make(chan ams.Response, 1)

	// register the handler.
	c.mu.Lock()
	if c.handler == nil {
		c.handler = make(map[uint32]chan ams.Response)
	}
	c.handler[pkt.Header().InvokeID] = h
	c.mu.Unlock()

	// send the request
	_, err := c.conn.Write(b.Bytes())
	if err != nil {
		c.mu.Lock()
		delete(c.handler, pkt.Header().InvokeID)
		c.mu.Unlock()
		return err
	}

	// wait for the response or timeout.
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(c.ReadTimeout):
		return ErrTimeout
	case r := <-h:
		return cb(r)
	}
}

// Read sends a Read request to the server.
func (c *Client) Read(ctx context.Context, r *ams.ReadRequest) (*ams.ReadResponse, error) {
	var resp *ams.ReadResponse
	err := c.send(ctx, r, func(r ams.Response) error {
		if x, ok := r.(*ams.ReadResponse); ok {
			resp = x
			return nil
		}
		return fmt.Errorf("got %T want %T", r, resp)
	})
	return resp, err
}

// ReadWrite sends a ReadWrite request to the server.
func (c *Client) ReadWrite(ctx context.Context, r *ams.ReadWriteRequest) (*ams.ReadWriteResponse, error) {
	var resp *ams.ReadWriteResponse
	err := c.send(ctx, r, func(r ams.Response) error {
		if x, ok := r.(*ams.ReadWriteResponse); ok {
			resp = x
			return nil
		}
		return fmt.Errorf("got %T want %T", r, resp)
	})
	return resp, err
}

// Write sends a Write request to the server.
func (c *Client) Write(ctx context.Context, r *ams.WriteRequest) (*ams.WriteResponse, error) {
	var resp *ams.WriteResponse
	err := c.send(ctx, r, func(r ams.Response) error {
		if x, ok := r.(*ams.WriteResponse); ok {
			resp = x
			return nil
		}
		return fmt.Errorf("got %T want %T", r, resp)
	})
	return resp, err
}

// GetSymHandleByName returns the offset of a variable.
func (c *Client) GetSymHandleByName(ctx context.Context, targetID, senderID ams.Addr, name string) (uint32, error) {
	req := ams.NewReadWriteRequest(targetID, senderID, ams.IdxGetSymHandleByName, 0, 4, []byte(name))
	res, err := c.ReadWrite(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("failed GetSymHandleByName %s: %s", name, err)
	}
	if res.Header().ErrorCode != ams.NoError {
		return 0, fmt.Errorf("failed ReadWrite: %d", res.Header().ErrorCode)
	}
	if len(res.Data) < 4 {
		return 0, fmt.Errorf("not enough data: %d", len(res.Data))
	}
	return binary.LittleEndian.Uint32(res.Data[:4]), nil
}
