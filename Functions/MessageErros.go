package ascii

import (
	"fmt"
	"os"
)

func MessageErrors() {
	fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]\n\nEX: go run . --output=<fileName.txt> something standard")
	os.Exit(0)
}
