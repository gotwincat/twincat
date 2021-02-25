// Copyright 2021 gotwincat authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ams

import (
	"fmt"
	"regexp"
	"strconv"
)

// reAddr is used for parsing a Twincat NetID.
var reAddr = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)\.(\d+)\.(\d+)\.(\d+):(\d+)`)

// Addr describes a Twincat NetID and port.
type Addr struct {
	NetID []byte
	Port  uint16
}

// String formats a Twincat NetID to a.b.c.d.e.f:port.
func (a Addr) String() string {
	return fmt.Sprintf("%d.%d.%d.%d.%d.%d:%d",
		a.NetID[0],
		a.NetID[1],
		a.NetID[2],
		a.NetID[3],
		a.NetID[4],
		a.NetID[5],
		a.Port,
	)
}

// MustParseAddr parses a NetID address and panics on error.
func MustParseAddr(s string) Addr {
	addr, err := ParseAddr(s)
	if err != nil {
		panic(err)
	}
	return addr
}

// ParseAddr parses a 'a.b.c.d.e.f:port' address into a NetID/Port
// for the AMS protocol. Values for a to e need must be between 0 and 255
// and the port must be between 0 and 65535.
func ParseAddr(s string) (Addr, error) {
	m := reAddr.FindStringSubmatch(s)
	if m == nil || len(m) != 8 {
		return Addr{}, fmt.Errorf("invalid address: %s", s)
	}

	netid := make([]byte, 6)
	for i := 1; i <= 6; i++ {
		n, err := strconv.ParseUint(m[i], 10, 32)
		if err != nil || n > 255 {
			return Addr{}, fmt.Errorf("invalid address: %s", s)
		}
		netid[i-1] = byte(n)
	}

	port, err := strconv.ParseUint(m[7], 10, 32)
	if err != nil || port > 65535 {
		return Addr{}, fmt.Errorf("invalid port: %s", s)
	}

	return Addr{netid, uint16(port)}, nil
}

func (a *Addr) Encode(b *Buffer) error {
	b.Write(a.NetID)
	b.WriteUint16(a.Port)
	return b.Err()
}

func (a *Addr) Decode(b *Buffer) error {
	a.NetID = b.ReadN(6)
	a.Port = b.ReadUint16()
	return b.Err()
}
