package tagger

import (
	"fmt"
	"path/filepath"
	"strings"
)

/*

line format

filename: tag,tag,tag

*/

const (
	fileToTagDivider  = ":"
	interTagDelimiter = ","
)

type taggedFileName string
type tagValue string

func deserialiseTags(l string) (taggedFileName, map[tagValue]bool) {
	m := map[tagValue]bool{}

	parts := strings.SplitN(l, fileToTagDivider, 2)
	if len(parts) == 1 {
		return taggedFileName(parts[0]), m
	}

	tags := strings.Split(parts[1], interTagDelimiter)
	for _, t := range tags {
		m[tagValue(t)] = true
	}

	return taggedFileName(parts[0]), m

}

func serialiseTags(filename taggedFileName, tags map[tagValue]bool) string {
	tagSl := make([]string, len(tags))
	i := 0
	for t := range tags {
		tagSl[i] = strings.Trim(string(t), " ")
		i++
	}

	tagStr := strings.Join(tagSl, interTagDelimiter)

	return fmt.Sprintf("%s%s%s", filename, fileToTagDivider, tagStr)
}

type TagFile map[taggedFileName]map[tagValue]bool

func unmarshalTagFile(b []byte) TagFile {
	s := string(b)
	lines := strings.Split(s, "\n")

	t := make(TagFile)

	for _, line := range lines {
		if line == "" {
			continue
		}
		fileName, lineTags := deserialiseTags(line)
		t[fileName] = lineTags
	}

	return t
}

func marshalTagFile(t TagFile) []byte {
	lines := make([]string, len(t))

	i := 0
	for filename, tags := range t {
		l := serialiseTags(filename, tags)
		lines[i] = l
		i++
	}
	s := strings.Join(lines, "\n")
	return []byte(s)
}

func validateTags(tags []string) error {
	for _, t := range tags {
		if strings.Contains(t, " ") {
			return fmt.Errorf("tag cannot contain a space: %s", t)
		}
		if strings.Contains(t, ",") {
			return fmt.Errorf("tag cannot contain commas: %s", t)
		}
	}
	return nil
}

const tagFileName = ".tags"

func getTagFilePathForDir(dir string) string {
	return filepath.Join(dir, tagFileName)
}
func isTagFile(fp string) bool {
	return filepath.Base(fp) == tagFileName
}
