package main

import "fmt"

func markdownLookup(path string) string {
	path = fmt.Sprintf("markdown/%s.md", path)
	return path
}
