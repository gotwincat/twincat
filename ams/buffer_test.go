// Copyright 2021 gotwincat authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ams

import (
	"testing"

	"github.com/pascaldekloe/goe/verify"
)

func TestBufferUint16(t *testing.T) {
	var bw Buffer
	bw.WriteUint16(0x1234)
	verify.Values(t, "err", bw.Err(), nil)
	verify.Values(t, "bytes", bw.Bytes(), []byte{0x34, 0x12})

	br := NewBuffer(bw.Bytes())
	n := br.ReadUint16()
	verify.Values(t, "err", br.Err(), nil)
	verify.Values(t, "data", n, uint16(0x1234))
}

func TestBufferUint32(t *testing.T) {
	var bw Buffer
	bw.WriteUint32(0x1234)
	verify.Values(t, "err", bw.Err(), nil)
	verify.Values(t, "bytes", bw.Bytes(), []byte{0x34, 0x12, 0x0, 0x0})

	br := NewBuffer(bw.Bytes())
	n := br.ReadUint32()
	verify.Values(t, "err", br.Err(), nil)
	verify.Values(t, "data", n, uint32(0x1234))
}

func TestBufferUint32Slice(t *testing.T) {
	var bw Buffer
	bw.WriteUint32Slice([]uint32{0xa, 0xb})
	verify.Values(t, "err", bw.Err(), nil)
	verify.Values(t, "bytes", bw.Bytes(), []byte{
		0xa, 0x0, 0x0, 0x0,
		0xb, 0x0, 0x0, 0x0,
	})

	br := NewBuffer(bw.Bytes())
	a := br.ReadUint32Slice(2)
	verify.Values(t, "err", br.Err(), nil)
	verify.Values(t, "data", a, []uint32{0xa, 0xb})
}

func TestBufferFloat32(t *testing.T) {
	var bw Buffer
	bw.WriteFloat32(1)
	verify.Values(t, "err", bw.Err(), nil)
	verify.Values(t, "bytes", bw.Bytes(), []byte{0x0, 0x0, 0x80, 0x3f})

	br := NewBuffer(bw.Bytes())
	n := br.ReadFloat32()
	verify.Values(t, "err", br.Err(), nil)
	verify.Values(t, "data", n, float32(1))
}

func TestBufferFloat32Slice(t *testing.T) {
	var bw Buffer
	bw.WriteFloat32Slice([]float32{1, 2})
	verify.Values(t, "err", bw.Err(), nil)
	verify.Values(t, "bytes", bw.Bytes(), []byte{
		0x0, 0x0, 0x80, 0x3f,
		0x0, 0x0, 0x0, 0x40,
	})

	br := NewBuffer(bw.Bytes())
	a := br.ReadFloat32Slice(2)
	verify.Values(t, "err", br.Err(), nil)
	verify.Values(t, "data", a, []float32{1, 2})
}
