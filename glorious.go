package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"

	"github.com/sstallion/go-hid"
)

const CFG_SIZE_USED = 131

type RGBEffect uint8

type DPI uint8

type GloriousMouse struct {
	vid  uint16
	pid  uint16
	name string
	path string
}

type RGB8 struct {
	R, G, B uint8
}

type GloriousConfig struct {
	ReportID  uint8
	CommandID uint8
	Unk1      uint8

	ConfigWrite uint8
	Unk2        [6]uint8

	Config1   uint8
	DpiCount  uint8
	ActiveDpi uint8

	DpiEnabled uint8

	Dpi       [16]uint8
	DpiColour [8]RGB8

	RgbEffect RGBEffect

	GloriousMode      uint8
	GloriousDirection uint8

	SingleMode   uint8
	SingleColour RGB8

	Breathing7Mode        uint8
	Breathing7ColourCount uint8
	Breathing7Colours     [7]RGB8

	TailMode uint8
	Unk4     [33]uint8

	RaveMode    uint8
	RaveColours [2]RGB8

	WaveMode uint8

	Breathing1Mode   uint8
	Breathing1Colour RGB8

	Unk5 uint8

	LiftOffDistance uint8
}

var SupportedDevices = []*GloriousMouse{
	{vid: 0x258a, pid: 0x27, name: "Dream Machines DM5"},
	{vid: 0x258a, pid: 0x33, name: "Glorious Model D"},
	{vid: 0x258a, pid: 0x36, name: "Glorious Model O/O-"},
}

//RGB Effect Enums
const (
	RGB_OFF        RGBEffect = 0
	RGB_GLORIOUS   RGBEffect = 0x1
	RGB_SINGLE     RGBEffect = 0x2
	RGB_BREATHING  RGBEffect = 0x5
	RGB_BREATHING7 RGBEffect = 0x3
	RGB_BREATHING1 RGBEffect = 0xa
	RGB_TAIL       RGBEffect = 0x4
	RGB_RAVE       RGBEffect = 0x7
	RGB_WAVE       RGBEffect = 0x9
)

func (r RGBEffect) Name() (string, bool) {
	switch r {
	case RGB_OFF:
		return "Off", true
	case RGB_GLORIOUS:
		return "Glorious Mode", true
	case RGB_SINGLE:
		return "Single Color", true
	case RGB_BREATHING:
		return "RGB Breathing", true
	case RGB_BREATHING7:
		return "Seven-color Breathing", true
	case RGB_BREATHING1:
		return "Single color Breathing", true
	case RGB_TAIL:
		return "Tail Effect", true
	case RGB_RAVE:
		return "Two-color Rave", true
	case RGB_WAVE:
		return "Wave Effect", true
	default:
		return "", false
	}
}

func (c *GloriousConfig) Mode(effect RGBEffect) (uint8, bool) {
	switch effect {
	case RGB_GLORIOUS:
		return c.GloriousMode, true
	case RGB_SINGLE:
		return c.SingleMode, true
	case RGB_BREATHING7:
		return c.Breathing7Mode, true
	case RGB_BREATHING1:
		return c.Breathing1Mode, true
	case RGB_TAIL:
		return c.TailMode, true
	case RGB_RAVE:
		return c.RaveMode, true
	case RGB_WAVE:
		return c.WaveMode, true
	default:
		return 0x0, false
	}
}

func NameToRGBEffect(n string) (RGBEffect, bool) {
	switch n {
	case "Off":
		return RGB_OFF, true
	case "Single Color":
		return RGB_SINGLE, true
	case "RGB Breathing":
		return RGB_BREATHING, true
	case "Seven-color Breathing":
		return RGB_BREATHING7, true
	case "Single color Breathing":
		return RGB_BREATHING1, true
	case "Tail Effect":
		return RGB_TAIL, true
	case "Two-color Rave":
		return RGB_RAVE, true
	case "Wave Effect":
		return RGB_WAVE, true
	default:
		return RGB_OFF, false
	}
}

func (m *GloriousMouse) Path() (string, bool) {
	if m.path == "" {
		return "", false
	}
	return m.path, true
}

func (m *GloriousMouse) enumFunc(info *hid.DeviceInfo) error {
	if info.InterfaceNbr == 1 {
		m.path = info.Path
	}
	return nil
}

func GetDevice(m *GloriousMouse) {
	hid.Enumerate(m.vid, m.pid, m.enumFunc)
}

func dumpConfig() *GloriousConfig {
	return nil
}

// Commands start here

func GetConfig(dev *hid.Device) *GloriousConfig {
	conf := [6]byte{0x5, 0x11}
	res, err := dev.SendFeatureReport(conf[:])
	if err != nil || res != len(conf) {
		log.Fatalln("res:", res, "in get config cmd, go err:", err, "hid err:", dev.Error())
	}

	cfg := GloriousConfig{
		ReportID: 0x4,
	}

	emptyBytes := 520 - binary.Size(cfg)

	var buf bytes.Buffer
	err = binary.Write(&buf, binary.BigEndian, cfg)
	if err != nil {
		log.Fatalln(err)
	}
	err = binary.Write(&buf, binary.BigEndian, make([]byte, emptyBytes))
	if err != nil {
		log.Fatalln(err)
	}
	d := buf.Bytes()

	res, err = dev.GetFeatureReport(d)
	if err != nil || res < 1 {
		log.Fatalln("res:", res, "in read config, go err:", err, "hid err:", dev.Error())
	}
	fmt.Printf("Config read with %v bytes\n", res)

	return Read(d)
}

func SetConfig(dev *hid.Device, conf *GloriousConfig) error {
	conf.ConfigWrite = CFG_SIZE_USED - 8
	confBytes, err := conf.Write()
	if err != nil {
		return err
	}

	res, err := dev.SendFeatureReport(confBytes)
	if res == -1 || err != nil {
		return errors.New(fmt.Sprintf("error writing config: (%v) res: %v", err, res))
	}
	return nil
}

func GetFirmwareVersion(dev *hid.Device) string {

	version := [6]byte{0x5, 0x1}
	res, err := dev.SendFeatureReport(version[:])
	if err != nil || res != len(version) {
		log.Fatalln("res:", res, "in get firmware version cmd, go err:", err, "hid err:", dev.Error())
	}

	res, err = dev.GetFeatureReport(version[:])
	if err != nil || res != len(version) {
		log.Fatalln("res:", res, "in read firmware version, go err:", err, "hid err:", dev.Error())
	}

	return fmt.Sprintf("%s", version)
}

func GetDebounceTime(dev *hid.Device) int {

	debounce := [6]byte{0x5, 0x1a}
	res, err := dev.SendFeatureReport(debounce[:])
	if err != nil || res != len(debounce) {
		log.Fatalln("res:", res, "in get debounce time cmd, go err:", err, "hid err:", dev.Error())
	}

	res, err = dev.GetFeatureReport(debounce[:])
	if err != nil || res != len(debounce) {
		log.Fatalln("res:", res, "in read debounce time, go err:", err, "hid err:", dev.Error())
	}

	return int(debounce[2] * 2)
}

func SetDebounceTime(dev *hid.Device, dbt int) error {

	debounce := [6]byte{0x5, 0x1a, byte(dbt / 2)}
	res, err := dev.SendFeatureReport(debounce[:])
	if err != nil || res != len(debounce) {
		log.Fatalln("res:", res, "in set debounce time, go err:", err, "hid err:", dev.Error())
	}
	return err
}

func (c *GloriousConfig) SetActiveDPI(opt int) error {
	if opt < 1 || opt > 6 {
		return errors.New("opt too high or low")
	}
	c.ActiveDpi = uint8(opt)

	return nil
}

func (c *GloriousConfig) SetDPI(opt int, dpi int) error {
	opt--
	if dpi < 200 {
		return errors.New("dpi is too low")
	}
	if opt < 0 || opt > 5 {
		return errors.New("opt too high or low")
	}

	c.Dpi[opt] = uint8(dpi/100 - 1)

	return nil
}

func (c *GloriousConfig) GetRGBMode() (int, int, error) {
	mode, ok := c.Mode(c.RgbEffect)
	if !ok {
		return 0, 0, errors.New("no fitting rgb mode")
	}
	brightness := int(mode >> 4)
	speed := int(mode & 0x0f)

	return brightness, speed, nil

}

func (c *GloriousConfig) SetRGBMode(brightness int, speed int) error {
	newModeSetting := uint8(speed) | uint8(brightness)<<4

	switch c.RgbEffect {
	case RGB_GLORIOUS:
		c.GloriousMode = newModeSetting
	case RGB_SINGLE:
		c.SingleMode = newModeSetting
	case RGB_BREATHING7:
		c.Breathing7Mode = newModeSetting
	case RGB_BREATHING1:
		c.Breathing1Mode = newModeSetting
	case RGB_TAIL:
		c.TailMode = newModeSetting
	case RGB_RAVE:
		c.RaveMode = newModeSetting
	case RGB_WAVE:
		c.WaveMode = newModeSetting
	}
	return nil

}

func (c *GloriousConfig) SetRGBEffect(effect RGBEffect) {
	c.RgbEffect = effect
}

func (c *GloriousConfig) SetRGBBrightness(brightness int) error {
	if brightness < 0 || brightness > 4 {
		return errors.New("brightness level too high or low")
	}

	switch c.RgbEffect {
	case RGB_GLORIOUS:
		return errors.New("brightness can not be set on RGB_GLORIOUS")
	}

	oldBrightness, speed, err := c.GetRGBMode()
	if err != nil {
		return err
	}

	fmt.Println("Old Brightness:", oldBrightness)

	fmt.Println("New Brightness:", brightness)

	return c.SetRGBMode(brightness, speed)

}

func (c *GloriousConfig) SetRGBSpeed(speed int) error {
	if speed < 0 || speed > 3 {
		return errors.New("brightness level too high or low")
	}

	brightness, oldSpeed, err := c.GetRGBMode()
	if err != nil {
		return err
	}

	fmt.Println("Old Speed:", oldSpeed)

	fmt.Println("New Speed:", speed)

	return c.SetRGBMode(brightness, speed)

}
