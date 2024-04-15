package stubs

import (
	"fmt"
	"os"
)

// ParseFile reads a file and returns its content as a byte slice.
func ParseFile(filename string) []byte {
	file, e := os.ReadFile(filename)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	return file
}
