package main

import (
	"fmt"
	"log"

	"github.com/sstallion/go-hid"
)

func main() {
	var mouse *GloriousMouse
	for _, mouse = range SupportedDevices {
		GetDevice(mouse)
		if _, ok := mouse.Path(); ok {
			break
		}
	}
	fmt.Printf("Device detected: %v (%v)\n", mouse.name, mouse.path)

	dev, err := hid.OpenPath(mouse.path)
	if err != nil {
		log.Fatalln(err)
	}
	dbt := GetDebounceTime(dev)
	fmt.Printf("Debounce time: %vms\n", dbt)

	fmt.Printf("Firmware version: %v", GetFirmwareVersion(dev))

	conf := GetConfig(dev)

	// TODO: Interactive / commandline based

	err = SetConfig(dev, conf)
	fmt.Println(err)
}
