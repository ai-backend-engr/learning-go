package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	header = `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Document</title>
		</head>
		<body>
	`

	footer = `	    
		</body>
		</html>
	`
)

func main() {
	// Parse flags
	fname := flag.String("file", "", "Markdown file to preview")
	flag.Parse()

	if *fname == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*fname); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(fname string) error {
	// Read all the data from the input file and check for errors
	data, err := os.ReadFile(fname)
	if err != nil {
		return err
	}

	htmlData := parseContent(data)

	outFname := fmt.Sprintf("%s.html", filepath.Base(fname))

	return saveHTML(outFname, htmlData)
}

func parseContent(input []byte) []byte {
	// Parse the markdown file through blackfriday and bluemonday to generate a valid and safe HTML
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	// Create a buffer of bytes to write to file
	var buffer bytes.Buffer

	// Write html to bytes buffer
	buffer.WriteString(header)
	buffer.Write(body)
	buffer.WriteString(footer)

	return buffer.Bytes()
}

func saveHTML(outFname string, data []byte) error {
	return os.WriteFile(outFname, data, 0644)
}
