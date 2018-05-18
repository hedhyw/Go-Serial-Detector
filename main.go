package main

import (
	"log"

	"github.com/hedhyw/Go-Serial-Detector/serialdet"
)

func main() {
	list, _ := serialdet.List()
	log.Print(list)
}
