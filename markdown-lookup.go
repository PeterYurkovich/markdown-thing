package main

import (
	"fmt"
	"io"
	"os"
)

func markdownLookup(path string) string {
	path = fmt.Sprintf("markdown/%s.md", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "no"
	}
	file, err := os.Open(path)
	if err != nil {
		return "no"
	}
	defer file.Close()
	buffer := make([]byte, 1024)
	output := ""
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "no"
		}
		output += string(buffer[:n])
	}
	return output
}
