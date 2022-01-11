package main

import (
	"fmt"
	"os"

	"github.com/hedhyw/Go-Serial-Detector/pkg/v1/serialdet"
)

func main() {
	list, err := serialdet.List()
	switch {
	case err != nil:
		fmt.Println(err)
		os.Exit(1)
	case len(list) == 0:
		fmt.Println("no serial ports found")
		os.Exit(0)
	}

	for i, port := range list {
		fmt.Printf("Index: %d\nDescription: %s\nPath: %s\n\n",
			i,
			port.Description(),
			port.Path(),
		)
	}
}
