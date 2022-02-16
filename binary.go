package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

type WriteGloriousConfig struct {
	ReportID  uint8
	CommandID uint8
	Unk1      uint8

	ConfigWrite uint8
	Unk2        [6]uint8

	Config1           uint8
	DpiCountActiveDpi uint8

	DpiEnabled uint8

	Dpi       [16]uint8
	DpiColour [8]RGB8

	RgbEffect uint8

	GloriousMode      uint8
	GloriousDirection uint8

	SingleMode   uint8
	SingleColour RGB8

	Breathing7Mode      uint8
	Breaing7ColourCount uint8
	Breathing7Colours   [7]RGB8

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

func Read(d []byte) *GloriousConfig {
	unk2Raw := d[4:]
	var unk2 [6]uint8
	copy(unk2[:], unk2Raw)

	dpiRaw := d[13:]
	var dpi [16]uint8
	copy(dpi[:], dpiRaw)

	dpiColourRaw := d[29:]
	var dpiColour [8]RGB8
	reader := bytes.NewReader(dpiColourRaw)
	binary.Read(reader, binary.BigEndian, &dpiColour)

	singleColourRaw := d[57:]
	var singleColour RGB8
	reader = bytes.NewReader(singleColourRaw)
	binary.Read(reader, binary.BigEndian, &singleColour)

	b7cRaw := d[62:]
	var b7c [7]RGB8
	reader = bytes.NewReader(b7cRaw)
	binary.Read(reader, binary.BigEndian, &b7c)

	unk4Raw := d[84:]
	var unk4 [33]uint8
	copy(unk4[:], unk4Raw)

	raveColoursRaw := d[118:123]
	var raveColours [2]RGB8
	reader = bytes.NewReader(raveColoursRaw)
	binary.Read(reader, binary.BigEndian, &raveColours)

	b1cRaw := d[126:]
	var b1c RGB8
	reader = bytes.NewReader(b1cRaw)
	binary.Read(reader, binary.BigEndian, &b1c)

	//raveColoursRaw := d[103:]
	//var raveColours [2]RGB8
	//reader = bytes.NewReader(raveColoursRaw)
	//binary.Read(reader, binary.BigEndian, &raveColours)

	cfg := &GloriousConfig{
		ReportID:              d[0],
		CommandID:             d[1],
		Unk1:                  d[2],
		ConfigWrite:           d[3],
		Unk2:                  unk2,
		Config1:               d[10],
		DpiCount:              d[11] & 4,
		ActiveDpi:             d[11] >> 4,
		DpiEnabled:            d[12],
		Dpi:                   dpi,
		DpiColour:             dpiColour,
		RgbEffect:             RGBEffect(d[53]),
		GloriousMode:          d[54],
		GloriousDirection:     d[55],
		SingleMode:            d[56],
		SingleColour:          singleColour,
		Breathing7Mode:        d[60],
		Breathing7ColourCount: d[61],
		Breathing7Colours:     b7c,
		TailMode:              d[83],
		Unk4:                  unk4,
		RaveMode:              d[117],
		RaveColours:           raveColours,
		WaveMode:              d[124],
		Breathing1Mode:        d[125],
		Breathing1Colour:      b1c,
		Unk5:                  d[130],
		LiftOffDistance:       d[131],
	}

	return cfg
}

func (c *GloriousConfig) Write() ([]byte, error) {
	// fuck the glorious firmware
	dpiCountActiveDpi := (c.ActiveDpi << 4) | c.DpiCount

	wc := &WriteGloriousConfig{
		c.ReportID, c.CommandID, c.Unk1, c.ConfigWrite, c.Unk2, c.Config1,
		dpiCountActiveDpi, c.DpiEnabled, c.Dpi, c.DpiColour, uint8(c.RgbEffect),
		c.GloriousMode, c.GloriousDirection, c.SingleMode, c.SingleColour,
		c.Breathing7Mode, c.Breathing7ColourCount, c.Breathing7Colours,
		c.TailMode, c.Unk4, c.RaveMode, c.RaveColours, c.WaveMode, c.Breathing1Mode,
		c.Breathing1Colour, c.Unk5, c.LiftOffDistance}

	var buf bytes.Buffer
	emptyBytes := 520 - binary.Size(wc)
	err := binary.Write(&buf, binary.BigEndian, wc)
	if err != nil {
		log.Fatalln(err)
	}
	err = binary.Write(&buf, binary.BigEndian, make([]byte, emptyBytes))
	if err != nil {
		log.Fatalln(err)
	}

	return buf.Bytes(), nil
}
