// Copyright 2021 gotwincat authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

// Package ams contains packet types for the Beckhoff Twincat3 AMS protocol.
//
// https://infosys.beckhoff.com/english.php?content=../content/1033/tc3_adsnetref/7312567947.html&id=
package ams

// Command ids. Order matters
const (
	CmdInvalid uint16 = iota
	CmdADSReadDeviceInfo
	CmdADSRead
	CmdADSWrite
	CmdADSReadState
	CmdADSWriteControl
	CmdADSAddDeviceNotification
	CmdADSDeleteDeviceNotification
	CmdADSDeviceNotification
	CmdADSReadWrite
)

// State Flags
const (
	StateResponse        = 1 << 0
	StateNoReturn        = 1 << 1
	StateADSCommand      = 1 << 2
	StateSystemCommand   = 1 << 3
	StateHighPrioCommand = 1 << 4
	StateTimestampAdded  = 1 << 5
	StateUDPCommand      = 1 << 6
	StateInitCommand     = 1 << 7
	StateBroadcast       = 1 << 15
)

const (
	NoError               = 0
	TargetMachineNotFound = 7
)

// IndexGroups
// https://infosys.beckhoff.com/english.php?content=../content/1033/tcadsdeviceplc/html/tcadsdeviceplc_indexadsservice.htm&id=
const (
	IdxGetSymHandleByName        = 0x0000F003
	IdxReserved                  = 0x0000F004
	IdxReadWriteSymValueByHandle = 0x0000F005
	IdxReleaseSymHandle          = 0x0000F006
	IdxReadIWriteI               = 0x0000F020
	IdxReadIXWriteIX             = 0x0000F021
	IdxADSIGRP_IOIMAGE_RISIZE    = 0x0000F025
	IdxReadQWriteQ               = 0x0000F030
	IdxReadQXWriteQX             = 0x0000F031
	IdxADSIGRP_IOIMAGE_ROSIZE    = 0x0000F035
	IdxADSIGRP_SUMUP_READ        = 0x0000F080
	IdxADSIGRP_SUMUP_WRITE       = 0x0000F081
	IdxADSIGRP_SUMUP_READWRITE   = 0x0000F082
)

// https://infosys.beckhoff.com/english.php?content=../content/1033/tc3_ads_intro/115845259.html&id=
const (
	PortAMSRouter            = 1
	PortTC3PLCRuntimeSystem1 = 851
)

// HasState returns true if the StateFlags in the header
// has the provided flags set.
func HasState(h AMSHeader, flag uint16) bool {
	return h.StateFlags&flag == flag
}

// Request is the interface for request objects.
type Request interface {
	Header() *AMSHeader
}

// Response is the interface for response objects.
type Response interface {
	Header() *AMSHeader
}

// Encoder is the interface for types that can be encoded
// to a buffer.
type Encoder interface {
	Encode(b *Buffer) error
}

// Decoder is the interface for types that can be decoded
// from a buffer.
type Decoder interface {
	Decode(b *Buffer) error
}
