// Copyright 2021 gotwincat authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ams

import (
	"reflect"
	"testing"

	"github.com/pascaldekloe/goe/verify"
)

var (
	target = MustParseAddr("1.2.3.4.5.6:1234")
	sender = MustParseAddr("5.6.7.8.9.0:5678")

	tcpHeader = TCPHeader{
		Reserved: 0x1234,
		Length:   0x12345678,
	}
	tcpHeaderBytes = []byte{
		0x34, 0x12, // Reserved
		0x78, 0x56, 0x34, 0x12, // Length
	}

	amsHeader = AMSHeader{
		Target:     target,
		Sender:     sender,
		CmdID:      0x1234,
		StateFlags: 0x5678,
		Length:     0x12345678,
		ErrorCode:  0x34567890,
		InvokeID:   0x56789012,
	}
	amsHeaderBytes = []byte{
		0x1, 0x2, 0x3, 0x4, 0x5, 0x6, // Target NetID
		0xd2, 0x04, // Target Port
		0x5, 0x6, 0x7, 0x8, 0x9, 0x0, // Sender NetID
		0x2e, 0x16, // Sender Port
		0x34, 0x12, // CmdID
		0x78, 0x56, // StateFlags
		0x78, 0x56, 0x34, 0x12, // Length
		0x90, 0x78, 0x56, 0x34, // ErrorCode
		0x12, 0x90, 0x78, 0x56, // InvokeID
	}
)

type codec interface {
	Encoder
	Decoder
}

func TestCodecTest(t *testing.T) {
	tests := []struct {
		name string
		p    codec
		b    []byte
	}{
		{
			name: "TCPHeader",
			p:    &tcpHeader,
			b:    tcpHeaderBytes,
		},
		{
			name: "AMSHeader",
			p:    &amsHeader,
			b:    amsHeaderBytes,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codecTest(t, tt.p, tt.b)
		})
	}
}

func codecTest(t *testing.T, c codec, b []byte) {
	t.Helper()

	var buf Buffer
	if err := c.Encode(&buf); err != nil {
		t.Fatalf("encode: %s", err)
	}

	verify.Values(t, "bytes", buf.Bytes(), b)

	cnew := reflect.New(reflect.TypeOf(c).Elem()).Interface().(codec)
	if err := cnew.Decode(NewBuffer(b)); err != nil {
		t.Fatalf("decode: %s", err)
	}

	verify.Values(t, "compare", cnew, c)
}
