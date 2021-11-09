package gendocs

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type JekyllSidebarNode struct {
	Title string
	URL   string
	F     []*JekyllSidebarNode
}

type JekyllSidebar struct {
	Name         string
	BasePagesUrl string

	rootNode       *JekyllSidebarNode
	nodesByPattern *JekyllSidebarNode
}

func NewJekyllSidebar(name, basePagesUrl string) *JekyllSidebar {
	return &JekyllSidebar{
		Name:         name,
		BasePagesUrl: basePagesUrl,
	}
}

func (s *JekyllSidebar) HandlePath(pathPattern string, _ []byte) error {
	markdownPagePath, err := PathPatternToFilesystemMarkdownPath(pathPattern)
	if err != nil {
		return err
	}

	if markdownPagePath == "index.md" {
		s.rootNode = &JekyllSidebarNode{
			Title: "Overview",
			URL:   fmt.Sprintf("/%s.html", path.Join(strings.TrimPrefix(s.BasePagesUrl, "/"), strings.TrimSuffix(markdownPagePath, ".md"))),
		}
	} else {
		if s.nodesByPattern == nil {
			s.nodesByPattern = &JekyllSidebarNode{
				Title: "Paths",
			}
		}

		s.nodesByPattern.F = append(s.nodesByPattern.F, &JekyllSidebarNode{
			Title: strings.TrimSuffix(markdownPagePath, ".md"),
			URL:   fmt.Sprintf("/%s.html", path.Join(strings.TrimPrefix(s.BasePagesUrl, "/"), strings.TrimSuffix(markdownPagePath, ".md"))),
		})
	}

	return nil
}

func (s *JekyllSidebar) getNodes() []*JekyllSidebarNode {
	return []*JekyllSidebarNode{s.rootNode, s.nodesByPattern}
}

func (s *JekyllSidebar) WriteFile(path string) error {
	var lines []string

	lines = append(lines, "# This file is generated by the github.com/werf/trdl/server/pkg/gendocs")
	lines = append(lines, "# DO NOT EDIT!")
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("%s: &%s", s.Name, s.Name))

	for _, node := range s.getNodes() {
		newLines, err := s.appendNode(lines, node, "  ")
		if err != nil {
			return err
		}
		lines = newLines
	}

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return fmt.Errorf("unable to mkdir %q: %s", filepath.Dir(path), err)
	}

	return os.WriteFile(path, append([]byte(strings.Join(lines, "\n")), '\n'), os.ModePerm)
}

func (s *JekyllSidebar) appendNode(lines []string, node *JekyllSidebarNode, indent string) ([]string, error) {
	lines = append(lines, fmt.Sprintf("%s- title: %s", indent, node.Title))

	if len(node.F) > 0 {
		lines = append(lines, fmt.Sprintf("%s  f:", indent))

		for _, subNode := range node.F {
			newLines, err := s.appendNode(lines, subNode, indent+"  ")
			if err != nil {
				return nil, err
			}
			lines = newLines
		}
	} else {
		lines = append(lines, fmt.Sprintf("%s  url: %s", indent, node.URL))
	}

	return lines, nil
}
