// Copyright 2021 gotwincat authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ams

// ReadWriteRequest is the packet for an AMS ReadWrite request.
type ReadWriteRequest struct {
	tcpHeader   TCPHeader
	amsHeader   AMSHeader
	IndexGroup  uint32
	IndexOffset uint32
	ReadLength  uint32
	WriteLength uint32
	Data        []byte
}

func NewReadWriteRequest(target, sender Addr, group, offset, readLength uint32, writeData []byte) *ReadWriteRequest {
	dataLen := uint32(len(writeData))
	return &ReadWriteRequest{
		tcpHeader: TCPHeader{
			Length: amsHeaderLen + dataLen + 16,
		},
		amsHeader: AMSHeader{
			Target:     target,
			Sender:     sender,
			CmdID:      CmdADSReadWrite,
			StateFlags: StateADSCommand,
			Length:     dataLen + 16,
		},
		IndexGroup:  group,
		IndexOffset: offset,
		ReadLength:  readLength,
		WriteLength: dataLen,
		Data:        writeData,
	}
}

func (r *ReadWriteRequest) Header() *AMSHeader {
	return &r.amsHeader
}

func (r *ReadWriteRequest) Encode(b *Buffer) error {
	b.WriteStruct(&r.tcpHeader)
	b.WriteStruct(&r.amsHeader)
	b.WriteUint32(r.IndexGroup)
	b.WriteUint32(r.IndexOffset)
	b.WriteUint32(r.ReadLength)
	b.WriteUint32(r.WriteLength)
	b.WriteN(r.Data, r.WriteLength)
	return b.Err()
}

func (r *ReadWriteRequest) Decode(b *Buffer) error {
	b.ReadStruct(&r.tcpHeader)
	b.ReadStruct(&r.amsHeader)
	r.IndexGroup = b.ReadUint32()
	r.IndexOffset = b.ReadUint32()
	r.ReadLength = b.ReadUint32()
	r.WriteLength = b.ReadUint32()
	r.Data = b.ReadN(int(r.WriteLength))
	return b.Err()
}

// ReadWriteResponse is the packet for an AMS ReadWrite response.
type ReadWriteResponse struct {
	tcpHeader TCPHeader
	amsHeader AMSHeader
	Result    uint32
	Length    uint32
	Data      []byte
}

func (r *ReadWriteResponse) Header() *AMSHeader {
	return &r.amsHeader
}

func (r *ReadWriteResponse) Encode(b *Buffer) error {
	b.WriteStruct(&r.tcpHeader)
	b.WriteStruct(&r.amsHeader)
	b.WriteUint32(r.Result)
	b.WriteUint32(r.Length)
	b.WriteN(r.Data, r.Length)
	return b.Err()
}

func (r *ReadWriteResponse) Decode(b *Buffer) error {
	b.ReadStruct(&r.tcpHeader)
	b.ReadStruct(&r.amsHeader)
	r.Result = b.ReadUint32()
	r.Length = b.ReadUint32()
	r.Data = b.ReadN(int(r.Length))
	return b.Err()
}

//
func IsReadWriteResponse(h AMSHeader) bool {
	return h.CmdID == CmdADSReadWrite && HasState(h, StateResponse)
}
