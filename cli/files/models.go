package files

import (
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	fileTreeNonTerminal string = "├──"
	fileTreeTerminal    string = "└──"
)

// Entry represents a file system entry.
type Entry struct {
	Name    string
	IsDir   bool
	Size    int64
	Mode    fs.FileMode
	ModTime time.Time
}

// TreePrinter is a struct that prints directory trees.
type TreePrinter struct {
	MaxDepth int // 0 = infinite
}

func NewTreePrinter(maxDepth int) *TreePrinter {
	return &TreePrinter{
		MaxDepth: maxDepth,
	}
}

// Render prints the given path as a tree view.
func (tp *TreePrinter) Render(path string) error {
	return tp.print(path, "", 0)
}

// internal recursive printer
func (tp *TreePrinter) print(path, prefix string, depth int) error {
	if tp.MaxDepth > 0 && depth >= tp.MaxDepth {
		return nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("cannot read %q: %w", path, err)
	}

	// Sort: dirs first, then files alphabetically
	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].IsDir() && !entries[j].IsDir() {
			return true
		}
		if !entries[i].IsDir() && entries[j].IsDir() {
			return false
		}
		return strings.ToLower(entries[i].Name()) < strings.ToLower(entries[j].Name())
	})

	for i, e := range entries {
		connector := "├──"
		if i == len(entries)-1 {
			connector = "└──"
		}

		name := e.Name()
		if e.IsDir() {
			name += "/"
		}

		fmt.Printf("%s%s %s\n", prefix, connector, name)

		if e.IsDir() {
			newPath := path + "/" + e.Name()
			newPrefix := prefix
			if i == len(entries)-1 {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}
			if err := tp.print(newPath, newPrefix, depth+1); err != nil {
				return err
			}
		}
	}

	return nil
}
