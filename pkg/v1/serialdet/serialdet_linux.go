// Package serialdet provides method for finding active serial ports
package serialdet

import (
	"bufio"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
)

const (
	udevSerialPath = "/dev/serial/by-id"
	sysfsTTYPath   = "/sys/class/tty/"
	sysfsUSBPrefix = "ttyUSB"
	sysfsDevUEvent = "device/uevent"
	devPath        = "/dev"
	rootID         = 0
)

var ueventDriverRe = regexp.MustCompile(`^.*DRIVER=(.+)$`)

func getProcFiles() [2]string {
	return [...]string{
		"/proc/tty/driver/serial",
		"/proc/tty/driver/usbserial",
	}
}

type listFunc func() ([]SerialPortInfo, error)

func getListFunctions() [3]listFunc {
	return [...]listFunc{
		procfsList, // for root user
		udevList,   // for regular user
		sysfsList,  // last hope : lists only /dev/ttyUSB*
	}
}

// isRoot checks that user has root privileges.
func isRoot() bool {
	cmd := exec.Command("id", "-u")

	out, err := cmd.Output()
	if err != nil {
		return false
	}

	userID, err := strconv.Atoi(strings.TrimSpace(string(out)))

	return err == nil && userID == rootID
}

// procfsList parses /proc/tty/driver/*serial.
func procfsList() ([]SerialPortInfo, error) {
	if !isRoot() {
		return nil, ErrPermissionDenied
	}

	var parser procfsParser
	ports := make([]SerialPortInfo, 0)
	for _, procFN := range getProcFiles() {
		parser.Reset()

		f, err := os.Open(procFN)
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			if err = parser.AddLine(scanner.Text()); err != nil {
				return nil, err
			}
		}

		if err = scanner.Err(); err != nil {
			return nil, err
		}

		ports = append(ports, parser.GetList()...)
	}

	return ports, nil
}

// uDevList parses /dev/serial/by-id.
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

// sysfsList returns only usb-serial devices
// using information from /sys/class/tty/*
func sysfsList() ([]SerialPortInfo, error) {
	files, err := ioutil.ReadDir(sysfsTTYPath)
	if err != nil {
		return nil, err
	}

	ports := make([]SerialPortInfo, 0)
	for _, file := range files {
		if strings.HasPrefix(file.Name(), sysfsUSBPrefix) {
			descr, err := getUeventInfo(path.Join(sysfsTTYPath, file.Name()))
			if err != nil {
				descr = file.Name()
			}

			info := SerialPortInfo{
				description: descr,
				path:        path.Join(devPath, file.Name()),
			}
			ports = append(ports, info)
		}
	}

	return ports, nil
}

// getUeventInfo returns a driver name of device
// by the information from /sys/class/tty/*/device/uevent
func getUeventInfo(p string) (string, error) {
	f, err := os.Open(path.Join(p, sysfsDevUEvent))
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		match := ueventDriverRe.FindStringSubmatch(scanner.Text())
		if len(match) == 2 {
			return match[1], nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", ErrDriverNotDefined
}

func list() (list []SerialPortInfo, err error) {
	for _, fun := range getListFunctions() {
		if list, err = fun(); err == nil {
			return list, nil
		}
	}

	return nil, err
}
