// Copyright 2021 gotwincat authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ams

// TCPHeader is the AMS/TCP packet header.
//
// https://infosys.beckhoff.com/english.php?content=../content/1033/tf7xxx_tc3_vision/27021602095817483.html&id=
type TCPHeader struct {
	Reserved uint16 // must be zero
	Length   uint32 // total length of header and data
}

func (h *TCPHeader) Encode(b *Buffer) error {
	b.WriteUint16(h.Reserved)
	b.WriteUint32(h.Length)
	return b.Err()
}

func (h *TCPHeader) Decode(b *Buffer) error {
	h.Reserved = b.ReadUint16()
	h.Length = b.ReadUint32()
	return b.Err()
}

const amsHeaderLen = 32 // 2*8 + 2 + 2 + 4 + 4 + 4

// AMSHeader is the AMS packet header.
//
// https://infosys.beckhoff.com/english.php?content=../content/1033/tf7xxx_tc3_vision/27021602095817483.html&id=
type AMSHeader struct {
	Target     Addr
	Sender     Addr
	CmdID      uint16
	StateFlags uint16
	Length     uint32
	ErrorCode  uint32
	InvokeID   uint32
}

func (h *AMSHeader) Encode(b *Buffer) error {
	b.WriteStruct(&h.Target)
	b.WriteStruct(&h.Sender)
	b.WriteUint16(h.CmdID)
	b.WriteUint16(h.StateFlags)
	b.WriteUint32(h.Length)
	b.WriteUint32(h.ErrorCode)
	b.WriteUint32(h.InvokeID)
	return b.Err()
}

func (h *AMSHeader) Decode(b *Buffer) error {
	b.ReadStruct(&h.Target)
	b.ReadStruct(&h.Sender)
	h.CmdID = b.ReadUint16()
	h.StateFlags = b.ReadUint16()
	h.Length = b.ReadUint32()
	h.ErrorCode = b.ReadUint32()
	h.InvokeID = b.ReadUint32()
	return b.Err()
}

// Header combines the AMS/TCP and AMS packet headers.
type Header struct {
	TCPHeader
	AMSHeader
}

func (r *Header) Encode(b *Buffer) error {
	b.WriteStruct(&r.TCPHeader)
	b.WriteStruct(&r.AMSHeader)
	return b.Err()
}

func (r *Header) Decode(b *Buffer) error {
	b.ReadStruct(&r.TCPHeader)
	b.ReadStruct(&r.AMSHeader)
	return b.Err()
}
