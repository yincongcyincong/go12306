package utils

import (
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