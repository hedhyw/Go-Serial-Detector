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
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

const udevSerialPath = "/dev/serial/by-id"
const sysfsTTYPath = "/sys/class/tty/"
const sysfsUSBPrefix = "ttyUSB"
const devPath = "/dev"
const rootID = 0

var procFiles = []string{
	"/proc/tty/driver/serial",
	"/proc/tty/driver/usbserial",
}

type listFunc func() ([]SerialPortInfo, error)

var listFunctions = []listFunc{
	procfsList,
	udevList,
	sysfsList,
}

func isRoot() bool {
	cmd := exec.Command("id", "-u")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	userID, err := strconv.Atoi(strings.TrimSpace(string(out)))
	return err == nil && userID == rootID
}

// procfsList parses /proc/tty/driver/*serial
func procfsList() ([]SerialPortInfo, error) {
	if !isRoot() {
		return nil, errors.New("Permission denied")
	}
	var parser procfsParser
	ports := make([]SerialPortInfo, 0)
	for _, procFN := range procFiles {
		parser.Reset()
		f, err := os.Open(procFN)
		if err != nil {
			return nil, err
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			parser.AddLine(scanner.Text())
		}
		if scanner.Err() != nil {
			return nil, err
		}
		ports = append(ports, parser.GetList()...)
	}
	return ports, nil
}

// uDevList parses /dev/serial/by-id
func udevList() ([]SerialPortInfo, error) {
	files, err := ioutil.ReadDir(udevSerialPath)
	if err != nil {
		return nil, err
	}
	ports := make([]SerialPortInfo, 0)
	for _, file := range files {
		fullPath := path.Join(udevSerialPath, file.Name())
		link, err := os.Readlink(fullPath)
		if err != nil {
			continue
		}
		absLinkPath := path.Join(udevSerialPath, link)
		info := SerialPortInfo{
			description: file.Name(),
			path:        absLinkPath,
		}
		ports = append(ports, info)
	}
	return ports, nil
}

// sysfsList returns only usb-serial devices without description
// using information from /sys/class/tty/*
func sysfsList() ([]SerialPortInfo, error) {
	files, err := ioutil.ReadDir(sysfsTTYPath)
	if err != nil {
		return nil, err
	}
	ports := make([]SerialPortInfo, 0)
	for _, file := range files {
		if strings.HasPrefix(file.Name(), sysfsUSBPrefix) {
			info := SerialPortInfo{
				description: file.Name(),
				path:        path.Join(devPath, file.Name()),
			}
			ports = append(ports, info)
		}
	}
	return ports, nil
}

func list() (list []SerialPortInfo, ok bool) {
	for _, fun := range listFunctions {
		if list, err := fun(); err == nil {
			return list, true
		}
	}
	return nil, false
}
