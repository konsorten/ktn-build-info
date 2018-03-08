package main

import (
	"fmt"
	"syscall"
)

func main() {
	winver, _ := syscall.GetVersion()

	fmt.Printf("Windows version: %v.%v.%v", byte(winver), byte(winver>>8), uint16(winver>>18))

	if byte(winver) == 6 && byte(winver>>8) == 2 {
		panic("Manifest is not active, as the default version 6.2 was returned (see https://github.com/golang/go/issues/17835)")
	}
}
