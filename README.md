# Go-Serial-Detector
This is a package that allows you to list active serial ports (/dev/ttyS*, /dev/ttyUSB*).

In linux systems information is obtained from udev, sysfs or procfs.

# OS support
This package currently supports only linux systems.

# Usage

```golang
import "github.com/hedhyw/Go-Serial-Detector/serialdet"
import "log"

if list, ok := serialdet.List(); ok {
  for _, p := range list {
    // p.Description() returns short information about serial port
    // p.Path() returns path to device, for example: "/dev/ttyUSB1"
    log.Print(p.Description(), " ", p.Path())
  }
}

```

# How to get
`go get github.com/hedhyw/Go-Serial-Detector/serialdet`