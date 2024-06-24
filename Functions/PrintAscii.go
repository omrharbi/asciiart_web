package ascii

import (
	"strings"
)

// Returning the ascci art of the text in argument
func PrintAsciArt(sl [][]string, s string) string {
	slice := [][]string{}
	result := ""
	for i := 0; i < len(s); i++ {
		if i+1 < len(s) && s[i] == '\\' && s[i+1] == 'n' {
			for j := 0; j < 8; j++ {
				for k := 0; k < len(slice); k++ {
					result += slice[k][j]
				}
				result += "\n"
				if len(slice) == 0 {
					break
				}
			}
			i++
			slice = [][]string{}
		} else {
			slice = append(slice, sl[rune(s[i]-32)])
		}
	}
	if len(slice) > 0 {
		for j := 0; j < 8; j++ {
			for k := 0; k < len(slice); k++ {
				result += slice[k][j]
			}
			result += "\n"
		}
	}
	if strings.HasSuffix(s, "\\n") {
		if len(s)-strings.Count(s, "\\n") != strings.Count(s, "\\n") {
			result += "\n"
		}
	}
	return result
}
