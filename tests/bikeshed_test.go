package tests

import (
	//	"/bikeshed"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNothing(t *testing.T) {
	//	t.Errorf("This test isn't actually written yet.")

	assert.Equal(t, "NotRight", "Right", "It's not right!")
}

var rgbToHsvTests = []struct {
	in  uint32
	out uint32
}{
	{ // black
		in:  0xFF000000,
		out: 0xFF000000,
	},
	{ // white
		in:  0xFFFFFFFF,
		out: 0xFFFFFFFF,
	},
	{ // nonsense
		in:  1,
		out: 2,
	},
}

func TestRgbToHsv(t *testing.T) {
//	t.Errorf("This test isn't implemented!")
	for _, test := range rgbToHsvTests {
		actual := Translate(test.in)
		assert.Equal(t, test.out, actual)
	}
}
