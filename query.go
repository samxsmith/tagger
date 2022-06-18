package tagger

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type QueryOpts struct {
	rootDir          string
	DisableRecursive bool
}

func QueryForFilenames(query string, opts *QueryOpts) ([]string, error) {
	// parse query
	// AND is implied with multiple tags
	// tag cannot have spaces as enforced by validateTags(...)
	// tags in query are split by space
	queryTags := strings.Split(query, " ")
	allMatches := make(TagFile)

	// find and read all tag files recursively
	filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
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

		matches := applyQuery(queryTags, tagFile)
		allMatches = appendToMap(allMatches, matches)

		return nil
	})

	matchingFileNames := make([]string, len(allMatches))

	i := 0
	for file := range allMatches {
		matchingFileNames[i] = string(file)
		i++
	}

	// perform query
	return matchingFileNames, nil
}

func applyQuery(queriedTags []string, tf TagFile) TagFile {
	for file, tags := range tf {
		if !fileHasAllQueriedTags(queriedTags, tags) {
			delete(tf, file)
		}
	}

	return tf
}

func fileHasAllQueriedTags(queriedTags []string, fileTags map[tagValue]bool) bool {
	for _, qt := range queriedTags {
		if isWildcard(qt) {
			if !fileTagsMatchWildcard(qt, fileTags) {
				return false
			}
			continue
		}
		if !fileTags[tagValue(qt)] {
			return false
		}
	}
	return true
}

func isWildcard(queryTag string) bool {
	return strings.HasSuffix(queryTag, "*")
}

func fileTagsMatchWildcard(queryTag string, fileTags map[tagValue]bool) bool {
	prefix := strings.TrimSuffix(queryTag, "*")
	for t := range fileTags {
		if strings.HasPrefix(string(t), prefix) {
			fmt.Println("WILDCARD MATCH", t)
			return true
		}
	}
	return false
}

func prependDirectoryToFileNames(dir string, tf TagFile) TagFile {
	dirTF := make(TagFile)
	for filename, v := range tf {
		fullPath := filepath.Join(dir, string(filename))
		dirTF[taggedFileName(fullPath)] = v
	}
	return dirTF
}

func appendToMap(m1, m2 TagFile) TagFile {
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}
