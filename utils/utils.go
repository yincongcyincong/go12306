package utils

import (
	"bytes"
	"image/png"
	"math/rand"
	"strings"
)

func GetBoolMap(strs []string) map[string]bool {
	res := make(map[string]bool)
	for _, s := range strs {
		res[s] = true
	}

	return res
}

func ReplaceSpecailChar(str string) string {
	str = strings.Replace(str, "%2A", "*", -1)
	str = strings.Replace(str, "%28", "(", -1)
	str = strings.Replace(str, "%29", ")", -1)

	return str
}

func GetRand(min, max int) int {
	if min >= max {
		return max
	}

	return rand.Intn(max-min) + min
}

func ReplaceChar(str string) string {
	str = strings.Replace(str, "%2A", "*", -1)
	str = strings.Replace(str, "%28", "(", -1)
	str = strings.Replace(str, "%29", ")", -1)
	str = strings.Replace(str, "%2F", "/", -1)
	str = strings.Replace(str, "+", "%20", -1)
	str = strings.Replace(str, "%3B", ";", -1)
	str = strings.Replace(str, "%2C", ",", -1)

	return str
}

func QrToString(img []byte) []string {

	imageByte := bytes.NewBuffer(img)
	image, err := png.Decode(imageByte)
	if err != nil {
		return []string{}
	}
	rectangle := image.Bounds()

	res := make([]string, rectangle.Max.Y)
	for i := rectangle.Min.Y; i < rectangle.Max.Y; i++ {
		for j := rectangle.Min.X; j < rectangle.Max.X; j++ {
			color := image.At(i, j)
			r, g, b, _ := color.RGBA()
			r = r >> 8
			g = g >> 8
			b = b >> 8
			if r == 255 && g == 255 && b == 255 {
				res[i] += "\u2588"
			} else {
				res[i] += " "
			}
		}
	}

	return res
}
