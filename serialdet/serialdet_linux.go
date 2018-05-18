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
	"io/ioutil"
	"os"
	"path"
)

const udevSerialPath = "/dev/serial/by-id"

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

func list() (list []SerialPortInfo, ok bool) {
	if list, err := udevList(); err == nil {
		return list, true
	}
	// @todo others
	return nil, false
}
