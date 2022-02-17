package ams

type ReadStateRequest struct {
	tcpHeader TCPHeader
	amsHeader AMSHeader
}

func NewReadStateRequest(target, sender Addr) *ReadStateRequest {
	return &ReadStateRequest{
		tcpHeader: TCPHeader{
			Length: amsHeaderLen,
		},
		amsHeader: AMSHeader{
			Target:     target,
			Sender:     sender,
			CmdID:      CmdADSReadState,
			StateFlags: StateADSCommand,
		},
	}
}

func (r *ReadStateRequest) Header() *AMSHeader {
	return &r.amsHeader
}

func (r *ReadStateRequest) Encode(b *Buffer) error {
	b.WriteStruct(&r.tcpHeader)
	b.WriteStruct(&r.amsHeader)
	return b.Err()
}

func (r *ReadStateRequest) Decode(b *Buffer) error {
	b.ReadStruct(&r.tcpHeader)
	b.ReadStruct(&r.amsHeader)
	return b.Err()
}

func IsReadStateRequest(h AMSHeader) bool {
	return h.CmdID == CmdADSReadState && h.StateFlags == StateADSCommand
}

type ReadStateResponse struct {
	tcpHeader   TCPHeader
	amsHeader   AMSHeader
	Result      uint32
	ADSState    uint16
	DeviceState uint16
}

func NewReadStateResponse(target, sender Addr, result uint32, adsState, deviceState uint16) *ReadStateResponse {
	return &ReadStateResponse{
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
		Result:      result,
		ADSState:    adsState,
		DeviceState: deviceState,
	}
}
func (r *ReadStateResponse) Header() *AMSHeader {
	return &r.amsHeader
}

func (r *ReadStateResponse) Encode(b *Buffer) error {
	b.WriteStruct(&r.tcpHeader)
	b.WriteStruct(&r.amsHeader)
	b.WriteUint32(r.Result)
	b.WriteUint16(r.ADSState)
	b.WriteUint16(r.DeviceState)
	return b.Err()
}

func (r *ReadStateResponse) Decode(b *Buffer) error {
	b.ReadStruct(&r.tcpHeader)
	b.ReadStruct(&r.amsHeader)
	r.Result = b.ReadUint32()
	r.ADSState = b.ReadUint16()
	r.DeviceState = b.ReadUint16()
	return b.Err()
}
