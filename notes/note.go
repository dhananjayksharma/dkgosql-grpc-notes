package notes

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Save a Note to the disk with the title as filename
func SaveToDisk(n *Note, folder string) error {
	filename := filepath.Join(folder, n.Title) //title should be sanitized
	filename = filename + n.Id
	fmt.Println("filename:", filename)
	return os.WriteFile(filename, n.Body, 0600)
}

// Scan files in a folder to find first occurrence of a keyword
func LoadFromDisk(keyword string, folder string) (*Note, error) {
	filename, err := searchKeywordInFilename(folder, keyword)
	if err != nil {
		return nil, err
	}
	body, err := os.ReadFile(filepath.Join(folder, filename))
	if err != nil {
		return nil, err
	}
	return &Note{Title: filename, Body: body}, nil
}

// Scan a directory and if a file name contains a substring, return the first one
func searchKeywordInFilename(folder string, keyword string) (string, error) {
	items, _ := ioutil.ReadDir(folder)
	for _, item := range items {

		// Read the whole file at once
		// this is the most inefficient search engine in the world
		// good enough for an example
		b, err := ioutil.ReadFile(filepath.Join(folder, item.Name()))
		if err != nil {
			// This is not normal but we can safely ignore it
			log.Printf("Could not read %v", item.Name())
		}
		s := string(b)

		if strings.Contains(s, keyword) {
			return item.Name(), nil
		}
	}
	return "", errors.New("no file contains this keyword")
}
