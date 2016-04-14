//go:directive to be ignored by gotomarkdown
/*
+++
title = "GoToMarkdown: Converting commented Go files to Markdown (plus some extras)"
date = "2016-04-14"
categories = ["tool"]
+++


*/

// # Imports and Globals

package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	commentPtrn      = `^\s*//\s?`
	commentStartPtrn = `^\s*/\*\s?`
	commentEndPtrn   = `\s?\*/\s*$`
	directivePtrn    = `^//go:`
	imagePtrn        = `!\[[^\]]+\(([^"\)]+)["\)]`
)

// imagePtrn should properly match:
// ![Alt text](path/to/img.jpg)
// ![Alt text](path/to/img.jpg "Title")
// ![Alt text](spaced path/to/some img.jpg)
// ![Alt text](spaced path/to/some img.jpg "Title")

var (
	comment          = regexp.MustCompile(commentPtrn)      // pattern for single-line comments
	commentStart     = regexp.MustCompile(commentStartPtrn) // pattern for /* comment delimiter
	commentEnd       = regexp.MustCompile(commentEndPtrn)   // pattern for */ comment delimiter
	directive        = regexp.MustCompile(directivePtrn)    // pattern for //go: directive, like //go:generate
	image            = regexp.MustCompile(imagePtrn)        // pattern for Markdown image tag
	allCommentDelims = regexp.MustCompile(commentPtrn + "|" + commentStartPtrn + "|" + commentEndPtrn)
	outDir           = flag.String("outdir", "out", "Output directory")
	dontCopyPics     = flag.Bool("nocopy", false, "Do not copy images to outdir")
)

// # First, some helper functions

// copyFiles copies a list of files or directories to a destination directory.
// The destination path must exist.
func copyFiles(dest string, files []string) (err error) {
	var result []byte
	for _, path := range files {
		result, err = exec.Command("copy", path, dest).Output()
		if err != nil {
			return errors.New(string(result) + "\n" + err.Error())
		}
	}
	return nil
}

// commentFinder returns a function that determines if the current line belongs to
// a comment region.
func commentFinder() func(string) bool {
	commentSectionInProgress := false
	return func(line string) bool {
		if comment.FindString(line) != "" {
			// "//" Comment line found.
			return true
		}
		// If the current line is at the start `/*` of a multi-line comment,
		// set a flag to remember we're within a multi-line comment.
		if commentStart.FindString(line) != "" {
			commentSectionInProgress = true
			return true
		}
		// At the end `*/` of a multi-line comment, clear the flag.
		if commentEnd.FindString(line) != "" {
			commentSectionInProgress = false
			return true
		}
		// The current line is within a `/*...*/` section.
		if commentSectionInProgress {
			return true
		}
		// Anything else is not a comment region.
		return false
	}
}

// isInComment returns true if the current line belongs to a comment region.
// A comment region `//` is either a comment line (starting with `//`) or
// a `/*...*/` multi-line comment.
var isInComment func(string) bool = commentFinder()

// isDirective returns true if the input argument is a Go directive,
// like `//go:generate`.
func isDirective(line string) bool {
	if directive.FindString(line) != "" {
		return true
	}
	return false
}

// extractMediaPath receives a line of text and searches for an image
// tag. If it finds one, it adds the path to the media list.
// NOTE: The function can only handle one image tag per line.
func extractMediaPath(line string) (path string, err error) {

	return path, nil
}

// convert receives a string containing commented Go code and converts it
// line by line into a Markdown document. Collect and return any media files
// found during this process.
func convert(in string) (out string, media []string, err error) {
	const (
		neither = iota
		comment
		code
	)
	lastLine := neither

	// Remove carriage returns.
	in = strings.Replace(in, "\r", "", -1)
	// Split at newline and process each line.
	for _, line := range strings.Split(in, "\n") {
		// Skip the line if it is a Go directive like //go:generate
		if isDirective(line) {
			continue
		}
		// Determine if the line belongs to a comment.
		if isInComment(line) {
			// Close the code block if a new comment begins.
			if lastLine == code {
				lastLine = comment
				out += "```\n"
			}
			// Detect `![image](path)` tags and add the path to the
			// media list.
			path, err := extractMediaPath(line)
			if err != nil {
				return "", nil, err
			}

			// Strip out any comment delimiter and add the line to the output.
			out += allCommentDelims.ReplaceAllString(line, "") + "\n"
		} else {
			// Open a new code block if the last line was a comment.
			if lastLine == comment {
				lastLine = code
				out += "```go"
			}
			// Add code lines verbatim to the output.
			out += line + "\n"
		}
	}
	return out, media, nil
}

// # Converting a file

// convertFile takes a file name, reads that file, converts it to
// Markdown, and writes it to `*outDir/&lt;basename>.md
func convertFile(filename string) (media []string, err error) {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Cannot read file " + filename + "\n" + err.Error())
	}
	name := filepath.Base(filename)
	ext := ".md"
	basename := name[:len(name)-3] // strip ".go"
	outname := filepath.Join(*outDir, basename) + ext
	md, media, err := convert(string(src))
	if err != nil {
		return nil, errors.New("Error converting " + filename + "\n" + err.Error())
	}
	err = os.MkdirAll(*outDir, 0744) // -rwxr--r--
	if err != nil {
		return nil, errors.New("Cannot create path: " + *outDir + " - Error: " + err.Error())
	}
	err = ioutil.WriteFile(outname, []byte(md), 0644) // -rw-r--r--
	if err != nil {
		return nil, errors.New("Cannot write file " + outname + " \n" + err.Error())
	}
	return media, nil
}

// # main

func main() {
	flag.Parse()
	for _, filename := range flag.Args() {
		log.Println("Converting", filename)
		media, err := convertFile(filename)
		if err != nil {
			log.Fatal("[Convert] " + err.Error())
		}
		if media != nil {
			log.Println("Copying media")
			err := copyFiles(*outDir, media)
			if err != nil {
				log.Fatal("[CopyFiles] " + err.Error())
			}
		}
	}
	log.Println("Done.")
}
