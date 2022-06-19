package tagger

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func ListTags(dir string, prefix string) ([]string, error) {
	combinedTagFile, err := digestTagFilesRecursively(dir)

	if err != nil {
		return nil, fmt.Errorf("error walking file tree: %w", err)
	}

	allTags := map[tagValue]bool{}
	for _, fileTags := range combinedTagFile {
		for t := range fileTags {
			allTags[t] = true
		}
	}

	if prefix != "" {
		allTagsFiltered := map[tagValue]bool{}
		for t := range allTags {
			if strings.HasPrefix(string(t), prefix) {
				allTagsFiltered[t] = true
			}
		}
		allTags = allTagsFiltered
	}

	var allTagsResult []string
	for t := range allTags {
		allTagsResult = append(allTagsResult, string(t))
	}

	return allTagsResult, nil
}

func digestTagFilesRecursively(root string) (TagFile, error) {
	combinedTagFile := make(TagFile)

	// find and read all tag files recursively
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !isTagFile(path) {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("open failed: %w", err)
		}
		defer f.Close()

		b, err := io.ReadAll(f)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		dir := filepath.Dir(path)

		tagFile := unmarshalTagFile(b)

		// prepend the directory so the result has
		// the correct path to the file
		tagFile = prependDirectoryToFileNames(dir, tagFile)
		combinedTagFile = appendToMap(combinedTagFile, tagFile)

		return nil
	})

	return combinedTagFile, err
}
