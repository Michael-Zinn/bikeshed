package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var colors = []struct {
	desc string
	rgb  uint32
	hsl  uint32
}{
	{
		desc: "black",
		rgb:  0xFF000000,
		hsl:  0xFF000000,
	},
	{
		desc: "red",
		rgb:  0xFFFF0000,
		hsl:  0xFF00FF7F,
	},
	{
		desc: "yellow",
		rgb:  0xFFFFFF00,
		hsl:  0xFF27FF7F,
	},
	{
		desc: "green",
		rgb:  0xFF00FF00,
		hsl:  0xFF55FF7F,
	},
	{
		desc: "cyan",
		rgb:  0xFF00FFFF,
		hsl:  0xFF7FFF7F,
	},
	{
		desc: "blue",
		rgb:  0xFF0000FF,
		hsl:  0xFFAAFF7F,
	},
	{
		desc: "magenta",
		rgb:  0xFFFF00FF,
		hsl:  0xFFD9FF7F,
	},
	
	{
		desc: "white",
		rgb:  0xFFFFFFFF,
		hsl:  0xFF0000FF,
	},
}

const hue uint32 = 0xFF0000
const sat uint32 = 0xFF00 // saturation
const lum uint32 = 0xFF   // luma

func TestRGBtoHSL(t *testing.T) {
	for _, test := range colors {
		hslFromCode := RGBtoHSL(test.rgb)
		assert.Equal(t, test.hsl&hue, hslFromCode&hue, test.desc+"\t: Hue")
		assert.Equal(t, test.hsl&sat, hslFromCode&sat, test.desc+"\t: Saturation")
		assert.Equal(t, test.hsl&lum, hslFromCode&lum, test.desc+"\t: Luminance")
//		assert.Equal(t, test.hsl, hslFromCode, test.desc+"\t: All")
	}
}

func TestHSLtoRGB(t *testing.T) {
	for _, test := range colors {
		rgbFromCode := HSLtoRGB(test.hsl)
		assert.Equal(t, test.rgb&hue, rgbFromCode&hue, test.desc+"\t: Hue")
		assert.Equal(t, test.rgb&sat, rgbFromCode&sat, test.desc+"\t: Saturation")
		assert.Equal(t, test.rgb&lum, rgbFromCode&lum, test.desc+"\t: Luminance")
//		assert.Equal(t, test.rgb, rgbFromCode, test.desc+"\t: All")
	}
}
