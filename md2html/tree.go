package md2html

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Tree struct {
	Name      string
	Link      string
	Directory bool
	Children  map[string]Tree
}

func (t Tree) GetNode(linkSegments []string) (Tree, error) {
	if len(linkSegments) == 1 {
		if linkSegments[0] == t.Link {
			return t, nil
		}
		return Tree{}, errors.New("no child")
	}
	child, ok := t.Children[linkSegments[0]]
	if !ok {
		return Tree{}, errors.New("no child")
	}
	return child.Children[linkSegments[0]].GetNode(linkSegments[1:])
}

func (t Tree) AddChild(linkSegments []string, fullLink string) error {
	if len(linkSegments) == 1 {
		t.Children[linkSegments[0]] = Tree{
			Name:      linkSegments[0],
			Link:      fullLink,
			Directory: false,
			Children:  map[string]Tree{},
		}
		return nil
	}
	child, ok := t.Children[linkSegments[0]]
	if !ok {
		child = Tree{
			Name:      linkSegments[0],
			Link:      "",
			Directory: true,
			Children:  map[string]Tree{},
		}
		t.Children[linkSegments[0]] = child
	}
	return child.AddChild(linkSegments[1:], fullLink)
}

func (t Tree) GetSortedChildren() []Tree {
	var children []Tree
	for _, child := range t.Children {
		children = append(children, child)
	}
	sort.Slice(children, func(i, j int) bool {
		// list directories first then sort both directories and files alphabetically
		if children[i].Directory && !children[j].Directory {
			return true
		}
		if !children[i].Directory && children[j].Directory {
			return false
		}
		return children[i].Name < children[j].Name
	})
	return children
}

func GetMarkdownTree() (Tree, error) {
	tree := Tree{
		Name:      "root",
		Link:      "/",
		Directory: true,
		Children:  map[string]Tree{},
	}
	err := filepath.WalkDir("markdown", func(path string, d os.DirEntry, err error) error {
		path = strings.TrimPrefix(path, "markdown/")
		if IgnoredLink(path) || BlockedLink(path) {
			return nil
		}
		if err != nil {
			return err
		}
		if d.IsDir() {
			// dont create the directory nodes until the children are added to prevent empty directories
			return nil
		}
		if strings.HasSuffix(d.Name(), ".md") {
			// Walk down the tree to the correct node, adding directories as needed
			if strings.Contains(path, "/") {
				tree.AddChild(strings.Split(path, "/"), path)
			} else {
				tree.AddChild([]string{path}, path)
			}
		}
		return nil
	})
	return tree, err
}

func GetLinkMarkdownMap() (map[string]string, error) {
	tree, err := GetMarkdownTree()
	if err != nil {
		return nil, err
	}
	linkMap := map[string]string{}
	for _, child := range tree.Children {
		markdown, err := MarkdownLookup(child.Link)
		if err != nil {
			return nil, err
		}
		linkMap[child.Link] = markdown
	}
	return linkMap, nil
}
