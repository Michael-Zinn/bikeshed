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

	L := (var_Max + var_Min) / 2
	//fmt.Println("LUMA: ", L)
	var H, S uint32

	if del_Max == 0 { //This is a gray, no chroma...
		//fmt.Println("NO CHROMA")
		H = 0 //HSL results from 0 to 1
		S = 0
	} else { //Chromatic data...
		if L < 128 {
			S = 255 * del_Max / (var_Max + var_Min)
		} else {
			S = 255 * del_Max / (511 - var_Max - var_Min)
		}

		del_R := (((var_Max - R) / 6) + (del_Max / 2)) / del_Max
		del_G := (((var_Max - G) / 6) + (del_Max / 2)) / del_Max
		del_B := (((var_Max - B) / 6) + (del_Max / 2)) / del_Max

		if R == var_Max {
			H = del_B - del_G
		} else if G == var_Max {
			H = (256 * 1 / 3) + del_R - del_B
		} else if B == var_Max {
			H = (256 * 2 / 3) + del_G - del_R
		}
		if H < 0 {
			H += 256
		}
		if H > 1 {
			H -= 256
		}
	}
	hsl = (rgb & 0xFF000000) | ((H & 0xFF) << 16) | ((S & 0xFF) << 8) | (L & 0xFF)
	//fmt.Println("RGB: ",rgb, toHex(rgb),"\t to HSL: ", hsl, toHex(hsl))
	return
}

func hueToRGB(v1, v2, vH uint32) uint32 { //Function Hue_2_RGB

	if vH < 0 {
		vH += 256
	}
	if vH > 1 {
		vH -= 256
	}

	if (6 * vH) < 255 {
		return (v1 + (v2-v1)*6*vH)
	}
	if (2 * vH) < 255 {
		return (v2)
	}
	if (3 * vH) < 511 {
		return (v1 + (v2-v1)*((511/3)-vH)*6)
	}
	return (v1)
}

func HSLtoRGB(hsl uint32) (rgb uint32) {

	H := (hsl >> 16) & 0xFF
	S := (hsl >> 8) & 0xFF
	L := hsl & 0xFF

	if S == 0 { //HSL from 0 to 1
		return 0xFF000000 | (L << 16) | (L << 8) | L
	} else {
		var var_2 uint32

		if L < 128 {
			var_2 = L * (1 + S)
		} else {
			var_2 = (L + S) - (S * L)
		}

		var_1 := 2*L - var_2

		R := hueToRGB(var_1, var_2, H+(1/3))
		G := hueToRGB(var_1, var_2, H)
		B := hueToRGB(var_1, var_2, H-(1/3))

		rgb = 0xFF000000 | (R << 16) | (G << 8) | B
		//fmt.Println("HSLtoRGB: hsl: ",toHex(hsl),"\trgb: ",toHex(rgb))
		return
	}

	// unreachable?
	return 0xDEADBEEF
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

	output, e := exec.Command("convert", imageFilePath /* "-scale", "1000x1000", */, "-gravity", "center" /*"-crop", "50%",*/, "-dither", "None", "-format", "%c", "histogram:info:").Output()

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
		//fmt.Println(saturationSums)
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
	var windowPixelCount uint32 = 0
	for _, pixelCount := range pixelCounts[:WINDOWSIZE] {
		//fmt.Println("Pixel counts: Hue: ",i,"\t Pixels for that hue: ", pixelCount)
		windowPixelCount += pixelCount
	}

	// set current window as initial max window
	maxWindowSaturationSum := windowSaturationSum
	maxWindowLeftPos := uint32(0)
	maxWindowPixelCount := windowPixelCount

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
			maxWindowPixelCount = windowPixelCount
		}
	}
	// maxSumLeftWindowPos should now point to the max saturation window

	// The final hue shall be the weighted average of the max saturation window
	// (The calculation is a bit confusing, the variable maxWindowConfusingHue
	// contains a value that doesn't map to any meaningful concept at all)
	var maxWindowConfusingHue uint32 = 0
	for windowPos := uint32(0); windowPos < WINDOWSIZE; windowPos++ {
		maxWindowConfusingHue += windowPos * saturationSums[(windowPos+maxWindowLeftPos)&255]
	}
	maxWindowAverageHue := (maxWindowConfusingHue/maxWindowPixelCount + maxWindowLeftPos) & 255

	// Convert back to ARGB for output
	//fmt.Println("LeftWindow: ", maxWindowLeftPos)
	//fmt.Println("Hue: ", maxWindowAverageHue)
	//fmt.Println("Saturation: ", averageSaturation)
	//fmt.Println("Luminance: ", averageBrightness)
	//fmt.Println("PixelCount: ", maxWindowPixelCount)

	//fmt.Println(maxWindowAverageHue, averageBrightness, averageSaturation)
	//	//fmt.Println(saturationSums)

	placeholderColor := (maxWindowAverageHue << 16) | (averageSaturation << 8) | averageBrightness
	fmt.Println("The HSL COLOR IS: ", 360*(placeholderColor>>16)/256, 100*((placeholderColor>>8)&0xFF)/256, 100*(placeholderColor&0xFF)/256)
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
	fmt.Println("THE PLACEHOLDER COLOR: ", toHex(color))

	fmt.Println("hsl test: f170a9", RGBtoHSL(0xf170a9))
}
