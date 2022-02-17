package ams

import (
	"testing"

	"github.com/pascaldekloe/goe/verify"
)

func TestNewReadStateResponse(t *testing.T) {
	got := NewReadStateResponse(target, sender, 0x1, 0x2, 0x3)
	want := &ReadStateResponse{
		tcpHeader: TCPHeader{
			Length: amsHeaderLen + 8,
		},
		amsHeader: AMSHeader{
			Target:     target,
			Sender:     sender,
			CmdID:      CmdADSReadState,
			StateFlags: StateADSCommand | StateResponse,
			Length:     8,
		},
		Result:      0x1,
		ADSState:    0x2,
		DeviceState: 0x3,
	}
	verify.Values(t, "", got, want)
}

func TestReadState(t *testing.T) {
	tests := []struct {
		name string
		p    codec
		b    []byte
	}{
		{
			name: "ReadStateRequest",
			p: &ReadStateRequest{
				tcpHeader: tcpHeader,
				amsHeader: amsHeader,
			},
			b: append(tcpHeaderBytes, amsHeaderBytes...),
		},
		{
			name: "ReadStateResponse",
			p: &ReadStateResponse{
				tcpHeader:   tcpHeader,
				amsHeader:   amsHeader,
				Result:      0x12345678,
				ADSState:    0x5678,
				DeviceState: 0x9012,
			},
			b: func() []byte {
				data := []byte{
					0x78, 0x56, 0x34, 0x12, // Result
					0x78, 0x56, // ADSState
					0x12, 0x90, // DeviceState
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
