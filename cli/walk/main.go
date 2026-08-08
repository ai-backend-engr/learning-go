package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type config struct {
	size int64
	ext  string
	list bool
}

func main() {
	// Parsing command line flags
	root := flag.String("root", ".", "Root directory to start")
	// Action flags
	list := flag.Bool("list", false, "List files only")
	// Filter flags
	ext := flag.String("ext", "", "File extension to filter out")
	size := flag.Int64("size", 0, "Minimum file size")

	cfg := config{size: *size, ext: *ext, list: *list}

	if err := run(*root, os.Stdout, cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(root string, out io.Writer, cfg config) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filterOut(path, cfg.ext, cfg.size, info) {
			return nil
		}

		// if list was explicitly set, dont do anything else
		if cfg.list {
			return listFile(path, out)
		}

		// List is the default option if nothing else was set
		return listFile(path, out)
	})
}
