package md2html

import (
	"slices"
	"strings"
)

func BlockedLink(link string) bool {
	block_list := []string{
		"Business/Accrisoft",
		"Business/Research",
		"Business/School",
		"Excalidraw",
		"Hobbies",
		"Personal",
		"Static",
	}
	return slices.ContainsFunc(block_list, func(block_item string) bool {
		return strings.HasPrefix(link, block_item)
	})
}

func IgnoredLink(link string) bool {
	ignore_list := []string{
		".git",
		".obsidian",
	}
	return slices.ContainsFunc(ignore_list, func(ignore_item string) bool {
		return strings.HasPrefix(link, ignore_item)
	})
}
