package main

// Warning: Written in a hurry, probably not idiomatic Go at all yet.

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// The hue histograms are 256 wide, WINDOWSIZE is the width
// of the hue histograms for which we want to find the max.
const WINDOWSIZE = 48

func RGBtoHSL(rgb uint32) (hsl uint32) {
	R := (rgb >> 16) & 0xFF
	G := (rgb >> 8) & 0xFF
	B := rgb & 0xFF

	var_Min := uint32(math.Min(float64(R), math.Min(float64(G), float64(B)))) //Min. value of RGB
	var_Max := uint32(math.Max(float64(R), math.Max(float64(G), float64(B)))) //Max. value of RGB
	del_Max := var_Max - var_Min                                              //Delta RGB value

	L := (var_Max + var_Min + 1) / 2
	//fmt.Println("LUMA: ", L)
	var H, S uint32

	if del_Max == 0 { //This is a gray, no chroma...
		//fmt.Println("NO CHROMA")
		H = 0 //HSL results from 0 to 1
		S = 0
	} else { //Chromatic data...
		if L < 128 {
			S = (255 * del_Max / (var_Max + var_Min)) * (256 - (L-128)*(L-128)/64) / 256
		} else {
			S = (255 * del_Max / (510 - var_Max - var_Min)) * (256 - (L-128)*(L-128)*(L-128)*(L-128)/1048576) / 256
		}

		//	fmt.Println(var_Max, del_Max)

		del_R := 255 * (((var_Max - R) / 6) + (del_Max / 2)) / del_Max
		del_G := 255 * (((var_Max - G) / 6) + (del_Max / 2)) / del_Max
		del_B := 255 * (((var_Max - B) / 6) + (del_Max / 2)) / del_Max

		if R == var_Max {
			H = del_B - del_G
		} else if G == var_Max {
			H = (255 * 1 / 3) + del_R - del_B
		} else if B == var_Max {
			H = (255 * 2 / 3) + del_G - del_R
		}
	}
	hsl = (rgb & 0xFF000000) | ((H & 0xFF) << 16) | ((S & 0xFF) << 8) | (L & 0xFF)
	return
}

func HSLtoRGB(hsl uint32) uint32 {

	h := float64((hsl>>16)&0xFF) / 255.0
	s := float64((hsl>>8)&0xFF) / 255.0
	l := float64(hsl&0xFF) / 255.0

	//fmt.Println("HSL", h, s, l)

	var r, g, b float64

	if s == 0 {
		ll := hsl & 0xFF
		return 0xFF000000 | (ll << 16) | (ll << 8) | ll // achromatic
	} else {
		var q float64
		if l < 0.5 {
			q = l * (1.0 + s)
		} else {
			q = l + s - l*s
		}
		p := 2.0*l - q
		//fmt.Println("qp", q, p)

		r = hue2rgb(p, q, h+1.0/3.0)
		g = hue2rgb(p, q, h)
		b = hue2rgb(p, q, h-1.0/3.0)
		//	fmt.Println("RGB", r, g, b)
	}

	return 0xFF000000 | (uint32(r*255) << 16) | (uint32(g*255) << 8) | uint32(b*255)
}

func hue2rgb(p, q, t float64) float64 {
	if t < 0 {
		t += 1.0
	}
	if t > 1 {
		t -= 1.0
	}

	if t < 1.0/6.0 {
		return p + (q-p)*6.0*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6.0
	}
	return p
}

// Uses ImageMagick to generate some basic histograms:
//   saturationSums: [hue byte]sum of saturations of pixels with that hue int
//   pixelCounts: [hue byte] number of pixels with that hue int
//   averageBrightness: average brightness of the whole image byte
//   averageSaturation: not sure yet, either of whole image or this value gets removed completely
func histograms(imageFilePath string) (saturationSums, pixelCounts [256]uint32, averageBrightness, averageSaturation uint32) {

	//fmt.Printf("Trying to read %s into ImageMagick\n", imageFilePath)

	//	output, e := exec.Command("convert o1.jpg -scale 100x100 -gravity center -crop 50% -dither None -remap color_map.gif -format %c histogram:info:"," ")).Output()

	//output, e := exec.Command("convert", "o1.jpg", "-scale", "100x100", "-gravity", "center", "-crop", "50%", "-dither", "None", "-remap", "color_map.gif", "-format", "%c", "histogram:info:").Output()

	output, e := exec.Command("convert", imageFilePath, "-scale", "128x128", "-gravity", "center" /*"-crop", "50%",*/, "-dither", "None", "-format", "%c", "histogram:info:").Output()

	//fmt.Println(string(output))

	//	output, e := exec.Command("convert", "o1.jpg", "-scale 100x100").Output()

	//"convert o1.jpg -scale 100x100 -gravity center -crop 50% -dither None -remap color_map.gif -format %c histogram:info:"

	//	output, e := exec.Command("convert", "o1.jpg", "-scale 100x100").Output()

	//	output, e := exec.Command("ls", "-l", "-a").Output()

	if e == nil {
		////fmt.Println("output: ", string(output))

		////fmt.Println(exec.Command("convert", "o1.jpg -scale 100x100 -gravity center -crop 50% -dither None -remap color_map.gif -format %c histogram:info: ").Start())

		// TODO load image file into histograms
		//	imagick.NewMagickWandFromImage

		saturationSums = [256]uint32{}
		pixelCounts = [256]uint32{}
		//for i := range saturationSums {
		//	pixelCounts[i] = rand.Intn(4000000)
		//	saturationSums[i] = rand.Intn(256) * pixelCounts[i]
		//}
		averageBrightness = 0
		averageSaturation = 0
		var pixelCount uint32 = 0

		lines := strings.Split(string(output), "\n")
		//lines look like this:
		//          6: (235, 26, 60) #EB1A3C rgb(235,26,60)
		//01234567890123456789012345689012345
		//0         11               29    35

		// //fmt.Println(lines)
		for _, line := range lines[0 : len(lines)-2] {
			//fmt.Println("Now trying to parse this: ",line)
			//			//fmt.Println(line[0:10],"----",line)
			countint64, _ := strconv.ParseInt(strings.TrimSpace(line[0:10]), 0, 32)
			count := (uint32)(countint64)
			//fmt.Println("count: ",count)
			//		//fmt.Println(count, line, counterror)
			rgbbytes, _ := hex.DecodeString(line[27:33])
			//fmt.Println("Read hex color: ",rgbbytes,"-----", line[29:35],"-----", line)
			var rgb uint32 = (uint32)(0xFF000000 | ((uint32)(rgbbytes[0]) << 16) | ((uint32)(rgbbytes[1]) << 8) | (uint32)(rgbbytes[2]))

			hsl := RGBtoHSL(rgb)
			//fmt.Println("hsl is ",hsl)
			h := (hsl >> 16) & 0xFF
			s := (hsl >> 8) & 0xFF
			l := hsl & 0xFF
			//fmt.Println(hsl,h,s,l)
			saturationSums[h] += s * count
			averageSaturation += s * count
			averageBrightness += l * count
			pixelCounts[h] += (uint32)(count)
			pixelCount += (uint32)(count)
		}
		//		fmt.Println(saturationSums)
		averageSaturation /= pixelCount
		averageBrightness /= pixelCount

		return

	} else {
		fmt.Println("error: ", e)
		fmt.Println("output: ", string(output))
		panic(e)
	}

	// unreachable code to satisfy the Go compiler?
	return
}

// imagefilepath should point to a file that can be handled by ImageMagick (nearly all pixel formats)
// returns the color as ARGB, A is currently always FF
//
// color is calculated in such a way that you can fade from a placeholder rectangle with that color
// to the actual image without irritating the user 
func placeholderColor(imageFilePath string) (color uint32) {

	saturationSums, pixelCounts, averageBrightness, averageSaturation := histograms(imageFilePath)
	//fmt.Println("Brightness: ", averageBrightness)

	// initialize window with beginning of histogram
	var windowSaturationSum uint32 = 0
	for _, saturationSum := range saturationSums[:WINDOWSIZE] {
		windowSaturationSum += saturationSum
	}
	//	fmt.Println("Initial Saturation Sum:", windowSaturationSum)

	var windowPixelCount uint32 = 0
	for _, pixelCount := range pixelCounts[:WINDOWSIZE] {
		//fmt.Println("Pixel counts: Hue: ",i,"\t Pixels for that hue: ", pixelCount)
		windowPixelCount += pixelCount
	}

	// set current window as initial max window
	maxWindowSaturationSum := windowSaturationSum
	maxWindowLeftPos := uint32(0)
//	maxWindowPixelCount := windowPixelCount

	// slide over the rest of the histogram to find the max saturation window
	for leftWindowPos := uint32(0); leftWindowPos < 256; leftWindowPos++ {
		// update running sum and pixel count (move 1 to the right)
		windowSaturationSum -= saturationSums[leftWindowPos]
		windowSaturationSum += saturationSums[(leftWindowPos+WINDOWSIZE)&255]
		windowPixelCount -= pixelCounts[leftWindowPos]
		windowPixelCount += pixelCounts[(leftWindowPos+WINDOWSIZE)&255]

		// check if that next window covers more saturation than the old one
		if windowSaturationSum > maxWindowSaturationSum {
			maxWindowSaturationSum = windowSaturationSum
			maxWindowLeftPos = leftWindowPos + 1
//			maxWindowPixelCount = windowPixelCount
			//			fmt.Println("Found better window:", maxWindowLeftPos, maxWindowPixelCount)
		}
	}
	// maxSumLeftWindowPos should now point to the max saturation window

	// The final hue shall be the weighted average of the max saturation window
	// (The calculation is a bit confusing, the variable maxWindowConfusingHue
	// contains a value that doesn't map to any meaningful concept at all)
	var maxWindowConfusingHue uint32 = 0
	for windowPos := uint32(0); windowPos < WINDOWSIZE; windowPos++ {
		maxWindowConfusingHue += windowPos * saturationSums[(windowPos+maxWindowLeftPos)&255]
		//		fmt.Println(maxWindowConfusingHue)
	}
	//	fmt.Println("confusing hue", maxWindowConfusingHue)
	maxWindowAverageHue := (uint32(float64(maxWindowConfusingHue)/float64(maxWindowSaturationSum)) + maxWindowLeftPos) & 255

	// Convert back to ARGB for output
	//	fmt.Println("LeftWindow: ", maxWindowLeftPos)
	//	fmt.Println("Hue: ", maxWindowAverageHue)
	//	fmt.Println("Saturation: ", averageSaturation)
	//	fmt.Println("Luminance: ", averageBrightness)
	//	fmt.Println("PixelCount: ", maxWindowPixelCount)

	//fmt.Println(maxWindowAverageHue, averageBrightness, averageSaturation)
	//	fmt.Println(saturationSums)

	//placeholderColor := (maxWindowAverageHue << 16) | (averageSaturation << 8) | averageBrightness
	//	fmt.Println("The HSL COLOR IS: ", 360*(placeholderColor>>16)/256, 100*((placeholderColor>>8)&0xFF)/256, 100*(placeholderColor&0xFF)/256)
	return HSLtoRGB(0xFF000000 | (maxWindowAverageHue << 16) | (averageSaturation << 8) | averageBrightness)
}

func toHex(color uint32) string {
	return hex.EncodeToString([]byte{byte((color >> 16) & 0xFF), byte((color >> 8) & 0xFF), byte(color & 0xFF)})
}

func main() {
	rand.Seed(time.Now().Unix())
	color := placeholderColor(os.Args[1])
	//	bytes := []byte{byte(color>>16&0xFF),byte(color>>8&0xFF),byte(color&0xFF)}
	////fmt.Print(hex.EncodeToString(bytes))
	fmt.Print(toHex(color))
}
