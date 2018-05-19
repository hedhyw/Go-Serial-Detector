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

// SerialPortInfo describes main information about active serial port
type SerialPortInfo struct {
	description string
	path        string
}

// Description contains serial ID or driver name of device
func (i SerialPortInfo) Description() string {
	return i.description
}

// Path is an absolute path to the device
func (i SerialPortInfo) Path() string {
	return i.path
}

// List returns active serial ports
func List() (ports []SerialPortInfo, ok bool) {
	return list()
}
