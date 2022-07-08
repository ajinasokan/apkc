package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// prompt prints the text and returns the string user entered. if empty then defaultValue is returned.
func prompt(text string, defaultValue string) string {
	fmt.Print(text)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			return defaultValue
		}
		return text
	}
	return defaultValue
}

// getFiles returns all the files in the path. if extension is not empty then that is also checked.
func getFiles(path string, extension string) []string {
	files := []string{}

	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		// ignoring traversal errors
		if err != nil {
			return nil
		}
		if f.IsDir() || strings.HasPrefix(f.Name(), ".") {
			return nil
		}
		if extension != "" {
			if strings.HasSuffix(f.Name(), extension) {
				files = append(files, path)
			}
		} else {
			files = append(files, path)
		}
		return nil
	})

	return files
}

// copyFiles copies files/directories. ref: https://stackoverflow.com/a/72246196
func copyFiles(src, dst string) error {
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// copy to this path
		outpath := filepath.Join(dst, strings.TrimPrefix(path, src))

		if info.IsDir() {
			os.MkdirAll(outpath, info.Mode())
			return nil // means recursive
		}

		// handle irregular files
		if !info.Mode().IsRegular() {
			switch info.Mode().Type() & os.ModeType {
			case os.ModeSymlink:
				link, err := os.Readlink(path)
				if err != nil {
					return err
				}
				return os.Symlink(link, outpath)
			}
			return nil
		}

		// copy contents of regular file efficiently

		// open input
		in, _ := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		// create output
		fh, err := os.Create(outpath)
		if err != nil {
			return err
		}
		defer fh.Close()

		// make it the same
		fh.Chmod(info.Mode())

		// copy content
		_, err = io.Copy(fh, in)
		return err
	})
}
