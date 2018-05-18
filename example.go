package main

import (
	"fmt"

	"github.com/hedhyw/Go-Serial-Detector/serialdet"
)

func main() {
	list, ok := serialdet.List()
	fmt.Println("Is it OK?", ok)
	for i, port := range list {
		fmt.Printf("\nIndex: %d\nDescription: %s\nPath: %s\n",
			i, port.Description(), port.Path())
	}
}
