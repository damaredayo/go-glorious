package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func runBefore(c *cli.Context) error {
	if c.Args().Len() == 0 {
		return fmt.Errorf("no arguments passed")
	}
	return nil
}

func (device *Device) runAfter(c *cli.Context) error {
	if c.Err() == nil {
		err := device.SetConfig()
		if err != nil {
			return fmt.Errorf("An error occured while writing the config: %v\n", err)
		}
		fmt.Println("Successfully updated configuration")
	}
	return nil
}
