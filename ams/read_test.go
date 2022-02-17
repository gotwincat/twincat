package ams

import (
	"testing"

	"github.com/pascaldekloe/goe/verify"
)

func TestNewReadRequest(t *testing.T) {
	got := NewReadRequest(target, sender, 0x1, 0x2, 0x3)
	want := &ReadRequest{
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
		IndexGroup:  0x1,
		IndexOffset: 0x2,
		Length:      0x3,
	}
	verify.Values(t, "", got, want)
}

func TestRead(t *testing.T) {
	tests := []struct {
		name string
		p    codec
		b    []byte
	}{
		{
			name: "ReadRequest",
			p: &ReadRequest{
				tcpHeader:   tcpHeader,
				amsHeader:   amsHeader,
				IndexGroup:  0x12345678,
				IndexOffset: 0x23456789,
				Length:      0x34567890,
			},
			b: func() []byte {
				data := []byte{
					0x78, 0x56, 0x34, 0x12, // IndexGroup
					0x89, 0x67, 0x45, 0x23, // IndexOffset
					0x90, 0x78, 0x56, 0x34, // Length
				}
				return append(append(tcpHeaderBytes, amsHeaderBytes...), data...)
			}(),
		},
		{
			name: "ReadResponse",
			p: &ReadResponse{
				tcpHeader: tcpHeader,
				amsHeader: amsHeader,
				Result:    0x12345678,
				Length:    0x3,
				Data:      []byte{0x00, 0x01, 0x02},
			},
			b: func() []byte {
				data := []byte{
					0x78, 0x56, 0x34, 0x12, // Result
					0x03, 0x00, 0x00, 0x00, // Length
					0x00, 0x01, 0x02, // Data
				}
				return append(append(tcpHeaderBytes, amsHeaderBytes...), data...)
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codecTest(t, tt.p, tt.b)
		})
	}
}
