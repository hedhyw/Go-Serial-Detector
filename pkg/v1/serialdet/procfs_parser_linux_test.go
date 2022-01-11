// Package serialdet provides method for finding active serial ports
package serialdet

import (
	"testing"
)

func TestProcfsParser(t *testing.T) {
	testCases := [...]struct {
		given []string
		want  []SerialPortInfo
	}{
		{
			given: []string{
				`usbserinfo:1.0 driver:2.0`,
				`0: module:ch341 name:"ch341-uart" vendor:1a86 product:7523 num_ports:1 port:0 path:usb-0000:00:14.0-2`,
			},
			want: []SerialPortInfo{
				{description: "ch341-uart", path: "/dev/ttyUSB0"},
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
				{description: "16550A", path: "/dev/ttyS0"},
				{description: "known", path: "/dev/ttyS27"},
			},
		},
	}

	var p procfsParser
	for i, c := range testCases {
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
			t.Fatalf(
				"Invalid GetList() length, given: %d, want: %d",
				len(list),
				len(c.want),
			)
		}

		for i, want := range c.want {
			if given := list[i]; given != want {
				t.Fatalf(
					"Invalid SerialPortInfo, given: %v, want: %v",
					given,
					want,
				)
			}
		}
	}
}
