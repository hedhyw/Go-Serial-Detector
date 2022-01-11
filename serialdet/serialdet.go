// Package serialdet is deprecated.
//
// Deprecated: Use the package github.com/hedhyw/Go-Serial-Detector/pkg/v1/serialdet instead.
package serialdet

import "github.com/hedhyw/Go-Serial-Detector/pkg/v1/serialdet"

// SerialPortInfo is deprecated.
//
// Deprecated: Use the package github.com/hedhyw/Go-Serial-Detector/pkg/v1/serialdet.
type SerialPortInfo = serialdet.SerialPortInfo

// List is deprecated.
//
// Deprecated: Use the package github.com/hedhyw/Go-Serial-Detector/pkg/v1/serialdet.
func List() (ports []SerialPortInfo, ok bool) {
	ports, err := serialdet.List()

	return ports, err == nil
}
