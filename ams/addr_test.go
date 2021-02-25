// Copyright 2021 gotwincat authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ams

import (
	"testing"

	"github.com/pascaldekloe/goe/verify"
)

func TestParseAddr(t *testing.T) {
	tests := []struct {
		s string
		a Addr
	}{
		{
			"1.2.3.4.5.6:1234",
			Addr{
				NetID: []byte{1, 2, 3, 4, 5, 6},
				Port:  1234,
			},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			a, err := ParseAddr(tt.s)
			if err != nil {
				t.Fatal(err)
			}
			verify.Values(t, "", a, tt.a)
		})
	}
}
