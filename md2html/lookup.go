package md2html

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var Error404 = errors.New(string(rune(http.StatusNotFound)))

func GetLinkFromName(name string) (linkInfo InternalLink, found bool) {
	// search the file tree in the markdown directory to see if there are one or more files with the same name
	// if there are, return the first one
	// if there are no files with the same name, return nil
	fileNames := []InternalLink{}
	err := filepath.WalkDir("markdown", func(path string, d os.DirEntry, err error) error {
		if IgnoredLink(path) {
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
			fileNames = append(fileNames, InternalLink{Link: strings.TrimPrefix(path, "markdown"), Name: strings.TrimSuffix(d.Name(), ".md")})
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

func MarkdownLookup(path string) (string, error) {
	if IgnoredLink(path) {
		return "", Error404
	}
	path = fmt.Sprintf("markdown/%s", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", Error404
	}
	file, err := os.Open(path)
	if err != nil {
		return "", err
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
			return "", err
		}
		output = append(output, buffer[:n]...)
	}
	return string(MdToHTML(output)), nil
}
