package ascii

import (
	"fmt"
	"os"
	"strings"
)

func WriteTextFileAscii(style string) ([][]string) {
	file, errs := os.ReadFile(style)
	if errs != nil {
		fmt.Println(errs)
		os.Exit(1)
	}
	started := ""
	sep := ""

	if style == "thinkertoy.txt" {
		started = string(file[2:])
		sep = string(file[:2])
	} else {
		started = string(file[1:])
		sep = string(file[:1])
	}
	splitFile := strings.Split(started, sep+sep)
	sliceToAppendAsci := [][]string{}
	for _, l := range splitFile {
		sliceToAppendAsci = append(sliceToAppendAsci, strings.Split(l, sep))
	}
	return sliceToAppendAsci
}
