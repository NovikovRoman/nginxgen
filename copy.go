package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func saveSnippets(snippetPath string, snippetName ...string) (err error) {
	for _, name := range snippetName {
		if err = copySnippet(snippetPath, name); err != nil && err != io.EOF {
			return
		}
	}
	return nil
}

func copySnippet(snippetPath string, name string) error {
	return copyFile(filepath.Join(templates, "snippets", name), filepath.Join(snippetPath, name), 1024)
}

func copyFile(src, dst string, buffersize int64) (err error) {
	var (
		sourceFileStat os.FileInfo
		destination    *os.File
		source         *os.File
	)

	if sourceFileStat, err = os.Stat(src); err != nil {
		return
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file. ", src)
	}

	if source, err = os.Open(src); err != nil {
		return err
	}
	defer func() {
		if derr := source.Close(); derr != nil {
			err = derr
		}
	}()

	if destination, err = os.Create(dst); err != nil {
		return err
	}
	defer func() {
		if derr := destination.Close(); derr != nil {
			err = derr
		}
	}()

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
