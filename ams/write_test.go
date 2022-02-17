package ams

import "testing"

func TestWrite(t *testing.T) {
	tests := []struct {
		name string
		p    codec
		b    []byte
	}{
		{
			name: "WriteRequest",
			p: &WriteRequest{
				tcpHeader:   tcpHeader,
				amsHeader:   amsHeader,
				IndexGroup:  0x12345678,
				IndexOffset: 0x23456789,
				Length:      0x3,
				Data:        []byte{0x00, 0x01, 0x02},
			},
			b: func() []byte {
				data := []byte{
					0x78, 0x56, 0x34, 0x12, // IndexGroup
					0x89, 0x67, 0x45, 0x23, // IndexOffset
					0x03, 0x00, 0x00, 0x00, // Length
					0x00, 0x01, 0x02, // Data
				}
				return append(append(tcpHeaderBytes, amsHeaderBytes...), data...)
			}(),
		},
		{
			name: "WriteResponse",
			p: &WriteResponse{
				tcpHeader: tcpHeader,
				amsHeader: amsHeader,
				Result:    0x12345678,
			},
			b: func() []byte {
				data := []byte{
					0x78, 0x56, 0x34, 0x12, // Result
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
