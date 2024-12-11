package html

import (
	"bytes"

	gohtml "golang.org/x/net/html"
)

type Node struct {
	*gohtml.Node
}

func Parse(b []byte) (*Node, error) {
	node, err := gohtml.Parse(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	return &Node{node}, nil
}

func (node *Node) Images() []*Node {
	if node.Type == gohtml.ElementNode && node.Data == "img" {
		return []*Node{node}
	}

	imgs := []*Node{}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		child := &Node{c}
		imgs = append(imgs, child.Images()...)
	}

	return imgs
}

func (node *Node) Attribute(key string) string {
	attr := ""
	for _, a := range node.Attr {
		if a.Key == key {
			attr = a.Val
		}
	}

	return attr
}

func (node *Node) Links() []string {
	links := []string{}
	if node.Type == gohtml.ElementNode && node.Data == "img" {
		return links
	}

	l := node.Attribute("href")
	if len(l) > 0 {
		links = append(links, l)
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		child := &Node{c}
		links = append(links, child.Links()...)
	}

	return links
}
