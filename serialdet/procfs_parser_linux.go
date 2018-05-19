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
	"errors"
	"regexp"
	"strings"
)

// HEADERS
const (
	procfsUSBHeadPrefix = "usbserinfo:1.0"
	procfsSerHeadPrefix = "serinfo:1.0"
)

// DEV PREFIXES
const (
	devUSBTTY = "/dev/ttyUSB"
	devSerTTY = "/dev/ttyS"
)

// PARSER REGEXPRS
const (
	procfsUSBRe = `^.+name:"(?P<Name>.+)" .+port:(?P<Port>\d+) .+$`
	procfsSerRe = `^(?P<Port>\d+):.+uart:(?P<Name>\w+) .+$`
)

const serialInvalidName = "unknown"

type serParserType byte
type parserInfo struct {
	re   *regexp.Regexp
	path string
}

var parserInfoByHead = map[string]parserInfo{
	procfsUSBHeadPrefix: parserInfo{
		re:   regexp.MustCompile(procfsUSBRe),
		path: devUSBTTY,
	},
	procfsSerHeadPrefix: parserInfo{
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
	// resolve type
	if p.info == nil {
		p.list = make([]SerialPortInfo, 0)
		for h, t := range parserInfoByHead {
			if strings.HasPrefix(line, h) {
				p.info = &t
				return nil
			}
		}
		return errors.New("Invalid information header")
	}

	var matched = p.info.re.FindStringSubmatch(line)
	if len(matched) != 3 {
		return errors.New("Invalid row")
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
