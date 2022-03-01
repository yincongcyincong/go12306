package utils

import (
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