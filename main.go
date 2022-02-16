package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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
				Subcommands: []*cli.Command{
					{
						Name:   "dpi",
						Usage:  "Set dpi",
						Before: runBefore,
						After:  device.runAfter,
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
						Before:  runBefore,
						After:   device.runAfter,
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
					{
						Name:    "lighting",
						Usage:   "Lighting configuration",
						Aliases: []string{"l"},
						Subcommands: []*cli.Command{
							{
								Name:    "effect",
								Aliases: []string{"e"},
								Before:  runBefore,
								After:   device.runAfter,
								Action: func(c *cli.Context) error {
									lightingName := strings.Join(c.Args().Slice(), " ")
									lightingMode, ok := NameToRGBEffect(lightingName)
									if !ok {
										return fmt.Errorf("invalid lighting mode")
									}
									device.conf.SetRGBEffect(lightingMode)
									return nil
								},
							},
							{
								Name:    "brightness",
								Aliases: []string{"b"},
								Before:  runBefore,
								After:   device.runAfter,
								Action: func(c *cli.Context) error {
									brightness, err := strconv.Atoi(c.Args().First())
									if err != nil {
										return err
									}
									return device.conf.SetRGBBrightness(brightness)
								},
							},
							{
								Name:    "speed",
								Aliases: []string{"s"},
								Before:  runBefore,
								After:   device.runAfter,
								Action: func(c *cli.Context) error {
									speed, err := strconv.Atoi(c.Args().First())
									if err != nil {
										return err
									}
									return device.conf.SetRGBSpeed(speed)
								},
							},
						},
					},
				},
			},
		},
		Version: "v1.1.0",
	}

	app.EnableBashCompletion = true
	err = app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
