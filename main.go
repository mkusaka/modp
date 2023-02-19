package main

import (
	"fmt"
	"io"
	"os"
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

	root, d := read(string(b))

	if err != nil {
		return err
	}

	putsDep(root, d, 1)
	return nil
}

func read(in string) (pkg, deps) {
	lines := strings.Split(in, "\n")

	dep := deps{}
	root := pkg{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		family := strings.SplitN(line, " ", 2)
		parent := newpkg(family[0])
		if parent.root() {
			root = parent
		}
		child := newpkg(family[1])
		if children, ok := dep[parent]; ok {
			dep[parent] = append(children, child)
		} else {
			dep[parent] = []pkg{child}
		}
	}

	return root, dep
}

func putsDep(p pkg, d deps, depth int) {
	indent := ""
	if depth >= 2 {
		indent = strings.Repeat("│ ", depth-1)
	}

	if p.root() {
		fmt.Printf("%s%s\n", indent, p.raw)
	}

	children := d[p]
	lastIndex := len(children) - 1
	for i, child := range children {
		grandchildren := d[child]
		if i == lastIndex {
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
			putsDep(child, d, depth+1)
		}
	}
}

type deps = map[pkg][]pkg

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
	name    string
	version string
	raw     string
}

func (p pkg) root() bool {
	return p.version == ""
}

func newpkg(nameversion string) pkg {
	name, version := splitNameVersion(nameversion)
	return pkg{
		name:    name,
		version: version,
		raw:     nameversion,
	}
}
