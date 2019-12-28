package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func saveSnippets(snippetPath string, snippetName ...string) error {
	for _, name := range snippetName {
		err := copySnippet(snippetPath, name)
		if err != nil && err != io.EOF {
			return err
		}
	}
	return nil
}

func copySnippet(snippetPath string, name string) error {
	return copyFile(filepath.Join(templates, "snippets", name), filepath.Join(snippetPath, name), 1024)
}

func copyFile(src, dst string, buffersize int64) (err error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file. ", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		if derr := source.Close(); derr != nil {
			err = derr
		}
	}()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if derr := destination.Close(); derr != nil {
			err = derr
		}
	}()

	if err != nil {
		return
	}

	buf := make([]byte, buffersize)
	n := 0
	for {
		n, err = source.Read(buf)
		if err != nil && err != io.EOF {
			return
		}
		if n == 0 {
			break
		}

		if _, err = destination.Write(buf[:n]); err != nil {
			return
		}
	}
	return
}
