// Copyright 2021 gotwincat authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ams

// ReadRequest is the packet for an AMS Read request.
type ReadRequest struct {
	tcpHeader   TCPHeader
	amsHeader   AMSHeader
	IndexGroup  uint32
	IndexOffset uint32
	Length      uint32
}

func NewReadRequest(target, sender Addr, group, offset, length uint32) *ReadRequest {
	return &ReadRequest{
		tcpHeader: TCPHeader{
			Length: amsHeaderLen + 12,
		},
		amsHeader: AMSHeader{
			Target:     target,
			Sender:     sender,
			CmdID:      CmdADSRead,
			StateFlags: StateADSCommand,
			Length:     12,
		},
		IndexGroup:  group,
		IndexOffset: offset,
		Length:      length,
	}
}

func (r *ReadRequest) Header() *AMSHeader {
	return &r.amsHeader
}

func (r *ReadRequest) Encode(b *Buffer) error {
	b.WriteStruct(&r.tcpHeader)
	b.WriteStruct(&r.amsHeader)
	b.WriteUint32(r.IndexGroup)
	b.WriteUint32(r.IndexOffset)
	b.WriteUint32(r.Length)
	return b.Err()
}

func (r *ReadRequest) Decode(b *Buffer) error {
	b.ReadStruct(&r.tcpHeader)
	b.ReadStruct(&r.amsHeader)
	r.IndexGroup = b.ReadUint32()
	r.IndexOffset = b.ReadUint32()
	r.Length = b.ReadUint32()
	return b.Err()
}

// ReadResponse is the packet for an AMS Read response.
type ReadResponse struct {
	tcpHeader TCPHeader
	amsHeader AMSHeader
	Result    uint32
	CBLength  uint32
	Data      []byte
}

func (r *ReadResponse) Header() *AMSHeader {
	return &r.amsHeader
}

func (r *ReadResponse) Encode(b *Buffer) error {
	b.WriteStruct(&r.tcpHeader)
	b.WriteStruct(&r.amsHeader)
	b.WriteUint32(r.Result)
	b.WriteUint32(r.CBLength)
	b.WriteN(r.Data, r.CBLength)
	return b.Err()
}

func (r *ReadResponse) Decode(b *Buffer) error {
	b.ReadStruct(&r.tcpHeader)
	b.ReadStruct(&r.amsHeader)
	r.Result = b.ReadUint32()
	r.CBLength = b.ReadUint32()
	r.Data = b.ReadN(int(r.CBLength))
	return b.Err()
}

// IsReadResponse returns true if the packet is a read response.
func IsReadResponse(h AMSHeader) bool {
	return h.CmdID == CmdADSRead && HasState(h, StateResponse)
}
