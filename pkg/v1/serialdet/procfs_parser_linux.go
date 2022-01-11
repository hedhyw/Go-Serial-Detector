// Package serialdet provides method for finding active serial ports
package serialdet

import (
	"regexp"
	"strings"
)

// Headers.
const (
	procfsUSBHeadPrefix = "usbserinfo:1.0"
	procfsSerHeadPrefix = "serinfo:1.0"
)

// Files to TTY.
const (
	devUSBTTY = "/dev/ttyUSB"
	devSerTTY = "/dev/ttyS"
)

// Regular expressions for parsing devices.
const (
	procfsUSBRe = `^.+name:"(?P<Name>.+)" .+port:(?P<Port>\d+) .+$`
	procfsSerRe = `^(?P<Port>\d+):.+uart:(?P<Name>\w+) .+$`
)

const serialInvalidName = "unknown"

type parserInfo struct {
	re   *regexp.Regexp
	path string
}

var parserInfoByHead = map[string]parserInfo{
	procfsUSBHeadPrefix: {
		re:   regexp.MustCompile(procfsUSBRe),
		path: devUSBTTY,
	},
	procfsSerHeadPrefix: {
		re:   regexp.MustCompile(procfsSerRe),
		path: devSerTTY,
	},
}

type procfsParser struct {
	info *parserInfo
	list []SerialPortInfo
}

func (p *procfsParser) Reset() {
	p.info = nil
}

func (p procfsParser) GetList() (res []SerialPortInfo) {
	res = make([]SerialPortInfo, len(p.list))
	copy(res, p.list)
	return res
}

func (p *procfsParser) AddLine(line string) error {
	// Resolve type
	if p.info == nil {
		p.list = make([]SerialPortInfo, 0)
		for h, t := range parserInfoByHead {
			if strings.HasPrefix(line, h) {
				p.info = &t

				return nil
			}
		}

		return ErrInvalidInformationHeader
	}

	var matched = p.info.re.FindStringSubmatch(line)
	if len(matched) != 3 {
		return ErrInvalidRow
	}

	var info SerialPortInfo
	for i, expName := range p.info.re.SubexpNames() {
		switch expName {
		case "Name":
			info.description = matched[i]
		case "Port":
			info.path = p.info.path + matched[i]
		}
	}

	if info.description != serialInvalidName {
		p.list = append(p.list, info)
	}

	return nil
}
