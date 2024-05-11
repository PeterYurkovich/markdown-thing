package md2html

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func GetLinkFromName(name string) (linkInfo InternalLink, found bool) {
	// search the file tree in the markdown directory to see if there are one or more files with the same name
	// if there are, return the first one
	// if there are no files with the same name, return nil
	fileNames := []InternalLink{}
	err := filepath.WalkDir("markdown", func(path string, d os.DirEntry, err error) error {
		if strings.Contains(path, ".git") || strings.Contains(path, ".obsidian") {
			return nil
		}
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(d.Name(), name+".md") {
			// remove a markdown/ prefix and a .md suffix
			fileNames = append(fileNames, InternalLink{Link: strings.TrimPrefix(path, "markdown/"), Name: strings.TrimSuffix(d.Name(), ".md")})
		}
		return nil
	})
	if err != nil {
		return InternalLink{}, false
	}
	if len(fileNames) == 0 {
		return InternalLink{}, false
	}
	return fileNames[0], true
}

func MarkdownLookup(path string) string {
	path = fmt.Sprintf("markdown/%s", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "no"
	}
	file, err := os.Open(path)
	if err != nil {
		return "no"
	}
	defer file.Close()
	buffer := make([]byte, 1024)
	output := []byte{}
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "no"
		}
		output = append(output, buffer[:n]...)
	}
	return string(MdToHTML(output))
}
