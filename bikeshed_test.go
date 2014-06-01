package main

import (
	"github.com/stretchr/testify/assert"
	"math"
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
//	{
//		desc: "dark green",
//		rgb:  0xFF001000,
//		hsl:  0xFF551E08,
//	},
	{
		desc: "red",
		rgb:  0xFFFF0000,
		hsl:  0xFF00FF80,
	},
	{
		desc: "yellow",
		rgb:  0xFFFFFF00,
		hsl:  0xFF2AFF80,
	},
	{
		desc: "green",
		rgb:  0xFF00FF00,
		hsl:  0xFF55FF80,
	},
	{
		desc: "cyan",
		rgb:  0xFF00FFFF,
		hsl:  0xFF7FFF80,
	},
	{
		desc: "blue",
		rgb:  0xFF0000FF,
		hsl:  0xFFAAFF80,
	},
	{
		desc: "magenta",
		rgb:  0xFFFF00FF,
		hsl:  0xFFD5FF80,
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

const maxError uint32 = 4

func TestRGBtoHSL(t *testing.T) {
	for _, test := range colors {
		hslFromCode := RGBtoHSL(test.rgb)
		//		assert.Equal(t, test.hsl&hue, hslFromCode&hue, test.desc+"\t: Hue")
		//	assert.Equal(t, test.hsl&sat, hslFromCode&sat, test.desc+"\t: Saturation")
		//assert.Equal(t, test.hsl&lum, hslFromCode&lum, test.desc+"\t: Luminance")
		//		assert.Equal(t, test.hsl, hslFromCode, test.desc+"\t: All")
		if channelDiff(test.hsl, hslFromCode) < maxError {
			t.Log(test.desc, "works")
		} else {
			//t.Error("Not good enough:",test.hsl,hslFromCode)
			// just for the pretty output
			assert.Equal(t, test.hsl, hslFromCode, test.desc+"\t: All")
		}
	}
}

func TestHSLtoRGB(t *testing.T) {
	for _, test := range colors {
		rgbFromCode := HSLtoRGB(test.hsl)

		//		assert.Equal(t, test.rgb, rgbFromCode, test.desc+"\t: All")

		if channelDiff(test.rgb, rgbFromCode) < maxError {
			t.Log(test.desc, "works")
		} else {
			//t.Error("Not good enough:",test.rgb,rgbFromCode)
			assert.Equal(t, test.rgb, rgbFromCode, test.desc+"\t: All")
		}
	}
}

func channelDiff(a, b uint32) uint32 {
	var diff float64 = 0

	diff += math.Abs(float64(a&0xFF) - float64(b&0xFF))
	diff += math.Abs(float64((a>>8)&0xFF) - float64((b>>8)&0xFF))
	diff += math.Abs(float64((a>>16)&0xFF) - float64((b>>16)&0xFF))

	return uint32(diff)
}
