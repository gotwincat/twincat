// Copyright 2021 gotwincat authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ams

// WriteRequest is the packet for an AMS write request.
type WriteRequest struct {
	tcpHeader   TCPHeader
	amsHeader   AMSHeader
	IndexGroup  uint32
	IndexOffset uint32
	Length      uint32
	Data        []byte
}

func NewWriteRequest(target, sender Addr, group, offset uint32, data []byte) *WriteRequest {
	dataLen := uint32(len(data))
	return &WriteRequest{
		tcpHeader: TCPHeader{
			Length: amsHeaderLen + dataLen + 12,
		},
		amsHeader: AMSHeader{
			Target:     target,
			Sender:     sender,
			CmdID:      CmdADSWrite,
			StateFlags: StateADSCommand,
			Length:     dataLen + 12,
		},
		IndexGroup:  group,
		IndexOffset: offset,
		Length:      dataLen,
		Data:        data,
	}
}

func (r *WriteRequest) Header() *AMSHeader {
	return &r.amsHeader
}

func (r *WriteRequest) Encode(b *Buffer) error {
	b.WriteStruct(&r.tcpHeader)
	b.WriteStruct(&r.amsHeader)
	b.WriteUint32(r.IndexGroup)
	b.WriteUint32(r.IndexOffset)
	b.WriteUint32(r.Length)
	b.Write(r.Data)
	return b.Err()
}

func (r *WriteRequest) Decode(b *Buffer) error {
	b.ReadStruct(&r.tcpHeader)
	b.ReadStruct(&r.amsHeader)
	r.IndexGroup = b.ReadUint32()
	r.IndexOffset = b.ReadUint32()
	r.Length = b.ReadUint32()
	r.Data = b.ReadN(int(r.Length))
	return b.Err()
}

// WriteResponse is the packet for an AMS write response.
type WriteResponse struct {
	tcpHeader TCPHeader
	amsHeader AMSHeader
	Result    uint32
}

func (r *WriteResponse) Header() *AMSHeader {
	return &r.amsHeader
}

func (r *WriteResponse) Encode(b *Buffer) error {
	b.WriteStruct(&r.tcpHeader)
	b.WriteStruct(&r.amsHeader)
	b.WriteUint32(r.Result)
	return b.Err()
}

func (r *WriteResponse) Decode(b *Buffer) error {
	b.ReadStruct(&r.tcpHeader)
	b.ReadStruct(&r.amsHeader)
	r.Result = b.ReadUint32()
	return b.Err()
}

// IsWriteResponse returns true if the packet is an AMS Write response.
func IsWriteResponse(h AMSHeader) bool {
	return h.CmdID == CmdADSWrite && HasState(h, StateResponse)
}
