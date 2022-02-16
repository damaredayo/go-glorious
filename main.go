package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/sstallion/go-hid"
	"github.com/urfave/cli/v2"
)

func main() {
	// Initalize mouse pointer then iterate over the supported devices until a match is found.
	var mouse *GloriousMouse
	for _, mouse = range SupportedDevices {
		GetDevice(mouse)
		if _, ok := mouse.Path(); ok {
			break
		}
	}

	if mouse == nil {
		fmt.Println("No supported mouse detected.")
		os.Exit(0)
	}

	fmt.Printf("Device detected: %v (%v)\n", mouse.name, mouse.path)

	dev, err := hid.OpenPath(mouse.path)
	if err != nil {
		fmt.Printf("%v, are you root?\n", err)
		os.Exit(1)
	}

	device := Device{dev, mouse, &GloriousConfig{}, 0}

	fmt.Printf("Firmware version: %v\n", device.GetFirmwareVersion())

	device.GetConfig()

	app := &cli.App{
		Name:  "go-glorious",
		Usage: "Configure Glorious devices on linux",
		Commands: []*cli.Command{
			{
				Name:    "set",
				Aliases: []string{"s"},
				Usage:   "set config option",
				Action: func(c *cli.Context) error {
					return fmt.Errorf("Please specify what to set")
				},
				Subcommands: []*cli.Command{
					{
						Name:  "dpi",
						Usage: "Set dpi",
						Action: func(c *cli.Context) error {
							dpiInt, err := strconv.Atoi(c.Args().First())
							if err != nil {
								return err
							}
							if err = device.conf.SetDPI(1, dpiInt); err != nil {
								return err
							}
							if err = device.conf.SetActiveDPI(1); err != nil {
								return err
							}
							fmt.Printf("Setting dpi to %vdpi\n", dpiInt)
							return nil
						},
					},
					{
						Name:    "debounce",
						Usage:   "Set debounce time",
						Aliases: []string{"db", "dbt"},
						Action: func(c *cli.Context) error {
							dbtInt, err := strconv.Atoi(c.Args().First())
							if err != nil {
								return err
							}
							if err = device.SetDebounceTime(dbtInt); err != nil {
								return err
							}
							fmt.Printf("Setting debounce time to %vms\n", dbtInt)
							return nil
						},
					},
				},
			},
		},
		Version: "v1.0.0",
	}

	app.EnableBashCompletion = true
	app.Run(os.Args)

	err = device.SetConfig()
	if err != nil {
		fmt.Printf("An error occured while writing the config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully updated configuration")
}
