// Package serialdet provides method for finding active serial ports
// Copyright 2018 Krivchun Maxim. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package serialdet

import (
	"testing"
)

var procfsParserCases = []struct {
	given []string
	want  []SerialPortInfo
}{
	{
		given: []string{
			`usbserinfo:1.0 driver:2.0`,
			`0: module:ch341 name:"ch341-uart" vendor:1a86 product:7523 num_ports:1 port:0 path:usb-0000:00:14.0-2`,
		},
		want: []SerialPortInfo{
			SerialPortInfo{description: "ch341-uart", path: "/dev/ttyUSB0"},
		},
	},
	{
		given: []string{
			`serinfo:1.0 driver revision:`,
			`0: uart:16550A port:000003F8 irq:4 tx:2725233 rx:2720703 brk:1 RTS|DTR`,
			`27: uart:known port:00000000 irq:0`,
			`28: uart:unknown port:00000000 irq:0`,
		},
		want: []SerialPortInfo{
			SerialPortInfo{description: "16550A", path: "/dev/ttyS0"},
			SerialPortInfo{description: "known", path: "/dev/ttyS27"},
		},
	},
}

func TestProcfsParser(t *testing.T) {
	var p procfsParser
	for i, c := range procfsParserCases {
		p.Reset()
		for _, in := range c.given {
			err := p.AddLine(in)
			if err != nil {
				t.Fatalf("error, case #%d: %v", i, err)
			}
		}
		list := p.GetList()
		if list == nil {
			t.Fatal("GetList() returns nil")
		}
		if len(list) != len(c.want) {
			t.Fatalf("Invalid GetList() length, given: %d, want: %d",
				len(list), len(c.want))
		}
		for i, want := range c.want {
			if given := list[i]; given != want {
				t.Fatalf("Invalid SerialPortInfo, given: %v, want: %v",
					given, want)
			}
		}
	}
}
