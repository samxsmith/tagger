package tagger

import (
	"fmt"
	"path/filepath"
	"strings"
)

type QueryOpts struct {
	rootDir          string
	DisableRecursive bool
}

func QueryForFilenames(dir string, queryTags []string, opts *QueryOpts) ([]string, error) {
	combinedTagFile, err := digestTagFilesRecursively(dir)

	if err != nil {
		return nil, fmt.Errorf("error walking file tree: %w", err)
	}

	matches := applyQuery(queryTags, combinedTagFile)
	matchingFileNames := make([]string, len(matches))

	i := 0
	for file := range matches {
		matchingFileNames[i] = string(file)
		i++
	}

	// perform query
	return matchingFileNames, nil
}

func applyQuery(queriedTags []string, tf TagFile) TagFile {
	res := make(TagFile)
	for file, tags := range tf {
		if fileHasAllQueriedTags(queriedTags, tags) {
			res[file] = tags
		}
	}

	return res
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
