// Package serialdet provides method for finding active serial ports
package serialdet

// SerialPortInfo describes main information about active serial port.
type SerialPortInfo struct {
	description string
	path        string
}

// Description contains serial ID or driver name of device.
func (i SerialPortInfo) Description() string {
	return i.description
}

// Path is an absolute path to the device.
func (i SerialPortInfo) Path() string {
	return i.path
}

// List returns active serial ports.
func List() (ports []SerialPortInfo, err error) {
	return list()
}
