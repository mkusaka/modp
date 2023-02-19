package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	root := read(string(b))

	putsDep(root, false, "", 1)
	return nil
}

type pkgMap = map[string]*pkg

func read(in string) *pkg {
	lines := strings.Split(in, "\n")

	sort.Sort(sort.Reverse(sort.StringSlice(lines)))
	pkgs := pkgMap{}

	rootRaw := &pkg{}

	for _, line := range lines {
		if line == "" {
			continue
		}
		family := strings.SplitN(line, " ", 2)
		parent := newpkg(family[0])
		if parent.root() {
			rootRaw = parent
		}
		if len(family) == 2 {
			child := newpkg(family[1])
			if cachedChild, ok := pkgs[child.raw]; ok {
				child = cachedChild
			}
			if p, ok := pkgs[parent.raw]; ok {
				p.children = append([]*pkg{child}, p.children...)
			} else {
				parent.children = append([]*pkg{child}, parent.children...)
				pkgs[parent.raw] = parent
			}

			if _, ok := pkgs[child.raw]; !ok {
				pkgs[child.raw] = child
			}
		} else {
			if _, ok := pkgs[parent.raw]; !ok {
				pkgs[parent.raw] = parent
			}
		}
	}

	return pkgs[rootRaw.raw]
}

func putsDep(p *pkg, eof bool, parentIndent string, depth int) {
	indent := parentIndent
	if depth >= 2 {
		if eof {
			indent += "  "
		} else {
			indent += "│ "
		}
	}

	if p.root() {
		fmt.Printf("%s%s\n", indent, p.raw)
	}

	children := p.children
	lastIndex := len(children) - 1
	for i, child := range children {
		grandchildren := child.children
		isLastIndex := i == lastIndex
		if isLastIndex {
			if len(grandchildren) >= 1 {
				fmt.Println(indent + "└─┬" + child.raw)
			} else {
				fmt.Println(indent + "└──" + child.raw)
			}
		} else {
			if len(grandchildren) >= 1 {
				fmt.Println(indent + "├─┬" + child.raw)
			} else {
				fmt.Println(indent + "├──" + child.raw)
			}
		}

		if len(grandchildren) >= 1 {
			putsDep(child, isLastIndex, indent, depth+1)
		}
	}
}
func splitNameVersion(str string) (name string, version string) {
	pair := strings.SplitN(str, "@", 2)
	if len(pair) == 2 {
		name = pair[0]
		version = pair[1]
		return name, version
	} else {
		name = pair[0]
		return name, ""
	}
}

type pkg struct {
	name     string
	version  string
	raw      string
	children []*pkg
}

func (p pkg) root() bool {
	return p.version == ""
}

func newpkg(nameversion string) *pkg {
	name, version := splitNameVersion(nameversion)
	return &pkg{
		name:    name,
		version: version,
		raw:     nameversion,
	}
}
