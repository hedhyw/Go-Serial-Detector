package serialdet

import (
	"errors"
	"regexp"
	"strings"
)

const usbserialInfo = "usbserinfo:1.0"
const serialInfo = "serinfo:1.0"
const usbTTY = "/dev/ttyUSB"
const comTTY = "/dev/ttyS"
const serialInvalidName = "unknown"

var sUSBRowRe = regexp.MustCompile(
	`^.+name:"(?P<Name>.+)" .+port:(?P<Port>\d+) .+$`)
var sCOMRowRe = regexp.MustCompile(
	`^(?P<Port>\d+):.+uart:(?P<Name>\w+) .+$`)

const (
	unknown   = 0
	usbserial = 1
	serial    = 2
)

type procfsParser struct {
	serialType int
	list       []SerialPortInfo
}

func (p *procfsParser) Reset() {
	p.serialType = unknown
}

func (p procfsParser) GetList() (res []SerialPortInfo) {
	res = make([]SerialPortInfo, len(p.list))
	copy(res, p.list)
	return res
}

func (p *procfsParser) AddLine(line string) error {
	if p.serialType == unknown {
		p.list = make([]SerialPortInfo, 0)
		if strings.HasPrefix(line, usbserialInfo) {
			p.serialType = usbserial
		} else if strings.HasPrefix(line, serialInfo) {
			p.serialType = serial
		} else {
			return errors.New("Invalid information head")
		}
		return nil
	}
	var matched []string
	var re *regexp.Regexp
	var path string
	if p.serialType == usbserial {
		re = sUSBRowRe
		path = usbTTY
	} else if p.serialType == serial {
		re = sCOMRowRe
		path = comTTY
	}
	matched = re.FindStringSubmatch(line)
	if len(matched) != 3 {
		return errors.New("Invalid row")
	}
	var info SerialPortInfo
	for i, expName := range re.SubexpNames() {
		switch expName {
		case "Name":
			info.description = matched[i]
		case "Port":
			info.path = path + matched[i]
		}
	}
	if info.description != serialInvalidName {
		p.list = append(p.list, info)
	}
	return nil
}
