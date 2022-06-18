package tagger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

const (
	logPrefix = "LOG:: "
)

func AddTag(filePath string, tags []string) error {
	if err := validateTags(tags); err != nil {
		return fmt.Errorf("invalid tag: %w", err)
	}
	// ensure passed file exists
	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("does that file exist? -> %w", err)
	}

	// find file dir
	dir, filename := filepath.Split(filePath)

	// ensure tag file
	f, err := createOrOpenTagFile(dir)
	if err != nil {
		return fmt.Errorf("error while checking config: %w", err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	// close file now we've read it
	f.Close()

	tagFile := unmarshalTagFile(b)

	if _, ok := tagFile[taggedFileName(filename)]; !ok {
		tagFile[taggedFileName(filename)] = map[tagValue]bool{}
	}

	for _, t := range tags {
		tagFile[taggedFileName(filename)][tagValue(t)] = true
	}

	tagFileB := marshalTagFile(tagFile)
	err = safeWriteTagFile(dir, tagFileB)

	return err
}

func safeWriteTagFile(dir string, b []byte) error {
	// first write new data to a temporary file
	// so that, if the write fails, we don't lose existing data

	// then move new file over existing
	// to overwrite with new data

	tmpName := fmt.Sprintf(".tags-%s", uuid.NewString())
	tmpPath := filepath.Join(dir, tmpName)
	if err := os.WriteFile(tmpPath, b, 0700); err != nil {
		return fmt.Errorf("failed to write updates: %w", err)
	}

	tagFilePath := getTagFilePathForDir(dir)
	if err := os.Rename(tmpPath, tagFilePath); err != nil {
		return fmt.Errorf("error while moving updates: %w", err)
	}

	return nil
}

func createOrOpenTagFile(dir string) (*os.File, error) {
	tagFilePath := getTagFilePathForDir(dir)

	// open, or create if doesn't exist
	f, err := os.OpenFile(tagFilePath, os.O_CREATE|os.O_RDWR, 0700)

	if err == nil {
		return f, nil
	}

	return nil, fmt.Errorf("failed to open tag file: %w", err)
}
