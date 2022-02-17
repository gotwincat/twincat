// Copyright 2021 gotwincat authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ams

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
)

// Buffer wraps a bytes.Buffer with an error and helper
// methods for encoding and decoding data structures
// in little endian encoding.
//
// All methods check if an error has occurred first
// and only return a meaningful value if there was no
// error so far.
type Buffer struct {
	b   bytes.Buffer
	err error
}

func NewBuffer(b []byte) *Buffer {
	buf := &Buffer{}
	buf.b.Write(b)
	return buf
}

// Bytes returns the content of the buffer.
func (buf *Buffer) Bytes() []byte {
	return buf.b.Bytes()
}

// Err returns the first error or nil.
func (buf *Buffer) Err() error {
	return buf.err
}

// Reset truncates the buffer to zero and resets the error.
func (buf *Buffer) Reset() {
	buf.b.Reset()
	buf.err = nil
}

// ReadFull reads len(b) bytes from the buffer.
// It returns io.ErrUnexpectedEOF if the buffer
// is too small.
func (buf *Buffer) ReadFull(b []byte) {
	if buf.err != nil {
		return
	}
	_, buf.err = io.ReadFull(&buf.b, b)
}

// ReadN reads n bytes from the buffer.
// It returns io.ErrUnexpectedEOF if the buffer
// is too small.
func (buf *Buffer) ReadN(n int) []byte {
	if buf.err != nil {
		return nil
	}
	b := make([]byte, n)
	_, buf.err = io.ReadFull(&buf.b, b)
	if buf.err != nil {
		return nil
	}
	return b
}

// ReadFloat32 reads a float32 from the buffer.
func (buf *Buffer) ReadFloat32() float32 {
	return math.Float32frombits(buf.ReadUint32())
}

// ReadFloat32Slice reads n float32 from the buffer.
func (buf *Buffer) ReadFloat32Slice(n int) []float32 {
	if buf.err != nil {
		return nil
	}
	a := buf.ReadUint32Slice(n)
	if buf.err != nil {
		return nil
	}
	fa := make([]float32, len(a))
	for i := range a {
		fa[i] = math.Float32frombits(a[i])
	}
	return fa
}

// ReadUint16 reads a uint16 from the buffer.
func (buf *Buffer) ReadUint16() uint16 {
	if buf.err != nil {
		return 0
	}
	return binary.LittleEndian.Uint16(buf.ReadN(2))
}

// ReadUint32 reads a uint32 from the buffer.
func (buf *Buffer) ReadUint32() uint32 {
	if buf.err != nil {
		return 0
	}
	return binary.LittleEndian.Uint32(buf.ReadN(4))
}

// ReadUint32Slice reads n uint32 from the buffer.
func (buf *Buffer) ReadUint32Slice(n int) []uint32 {
	if buf.err != nil {
		return nil
	}
	a := make([]uint32, n)
	for i := range a {
		a[i] = buf.ReadUint32()
	}
	return a
}

// ReadStruct decodes a type that implements the Decoder
// interface from the buffer.
func (buf *Buffer) ReadStruct(x Decoder) {
	if buf.err != nil {
		return
	}
	buf.err = x.Decode(buf)
}

// Write writes b to the buffer.
func (buf *Buffer) Write(b []byte) {
	if buf.err != nil {
		return
	}
	_, buf.err = buf.b.Write(b)
}

// WriteN writes up to n bytes of b to the buffer.
// This is similar to Write(b[:max]) for len(b) >= max.
func (buf *Buffer) WriteN(b []byte, max uint32) {
	if buf.err != nil {
		return
	}
	if len(b) > int(max) {
		b = b[:max]
	}
	buf.Write(b)
}

// WriteFloat32 writes a float32 to the buffer.
func (buf *Buffer) WriteFloat32(n float32) {
	buf.WriteUint32(math.Float32bits(n))
}

// WriteFloat32Slice writes all float32 without a length
// encoding to the buffer.
func (buf *Buffer) WriteFloat32Slice(a []float32) {
	aa := make([]uint32, len(a))
	for i := range a {
		aa[i] = math.Float32bits(a[i])
	}
	buf.WriteUint32Slice(aa)
}

// WriteUint16 writes a uint16 to the buffer.
func (buf *Buffer) WriteUint16(n uint16) {
	if buf.err != nil {
		return
	}
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, n)
	_, buf.err = buf.b.Write(b)
}

// WriteUint32 writes a uint32 to the buffer.
func (buf *Buffer) WriteUint32(n uint32) {
	if buf.err != nil {
		return
	}
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, n)
	_, buf.err = buf.b.Write(b)
}

// WriteUint32Slice writes all uint32 without a length
// encoding to the buffer.
func (buf *Buffer) WriteUint32Slice(a []uint32) {
	if buf.err != nil {
		return
	}
	b := make([]byte, 4*len(a))
	for i := range a {
		binary.LittleEndian.PutUint32(b[i*4:(i+1)*4], a[i])
	}
	_, buf.err = buf.b.Write(b)
}

// WriteStruct encodes a type that implements the
// Encoder interface to the buffer.
func (buf *Buffer) WriteStruct(x Encoder) {
	if buf.err != nil {
		return
	}
	buf.err = x.Encode(buf)
}
