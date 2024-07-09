package ascii

import (
	"strings"
	"sync"
)

func PrintAsciArt(sl [][]string, st  string) string {
	r := strings.Split(st, "\r\n")
	resultChan := make(chan string, len(r))
	var wg sync.WaitGroup

	for _, s := range r {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			slice := [][]string{}
			for i := 0; i < len(s); i++ {
				slice = append(slice, sl[rune(s[i]-32)])
			}

			lineResult := ""
			for j := 0; j < 8; j++ {
				for k := 0; k < len(slice); k++ {
					lineResult += slice[k][j]
				}
				lineResult += "\n"
			}
			resultChan <- lineResult
		}(s)
	}

	// Close the result channel once all goroutines are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	result := ""
	for lineResult := range resultChan {
		result += lineResult
	}

	return result
}
