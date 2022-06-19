package tagger

import (
	"fmt"
	"io"
	"os"
)

func Init(dir string) error {
	// if tagfile present, abort
	path := getTagFilePathForDir(dir)
	_, err := os.Stat(path)
	if err == nil {
		return fmt.Errorf(".tags file already exists here, can't re-initialise")
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("file system error: %w", err)
	}

	// confirmed: no existing tag file

	// list files -- non-recursive
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("error reading dir: %w", err)
	}

	f, err := createOrOpenTagFile(dir)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	tagFile := unmarshalTagFile(b)

	// for each, tag date
	for _, file := range files {
		if file.Name()[0] == '.' {
			// skip dot-files
			continue
		}
		fi, _ := file.Info()
		mt := fi.ModTime()
		dateTag := fmt.Sprintf("date:%s", mt.Format("2006-01-02"))

		tagFile = tagFile.AddTag(
			fi.Name(),
			[]string{dateTag},
		)
	}

	tagFileB := marshalTagFile(tagFile)
	err = safeWriteTagFile(dir, tagFileB)
	return err
}
