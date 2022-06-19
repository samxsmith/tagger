package e2e

import (
	"os"
	"sort"
	"testing"

	"github.com/samxsmith/tagger"
	"github.com/stretchr/testify/assert"
)

func TestInitTagFile(t *testing.T) {
	os.Remove("./test_data/init_data/.tags")

	// inits a file with dates
	err := tagger.Init("./test_data/init_data/")
	assert.Nil(t, err)

	// validate file
	b, err := os.ReadFile("./test_data/init_data/.tags")
	assert.Nil(t, err)

	s := string(b)
	expected := "file_one:date:2022-06-18\nfile_two:date:2022-06-18"
	assert.Equal(t, expected, s)

	// repeat, aborts
	err = tagger.Init("./test_data/init_data/")
	assert.NotNil(t, err)
}

func TestListTag(t *testing.T) {
	allTags, err := tagger.ListTags("./test_data/static_tag_files", "")
	assert.Nil(t, err)

	expected := []string{
		"person:grandad", "person:nan",
		"person:spouse", "occassion:wedding", "person:me",
		"occassion:holiday", "subject:sea", "subject:flower",
	}
	sort.Strings(allTags)
	sort.Strings(expected)

	assert.Equal(t, allTags, expected)
}

func TestRemoveTag(t *testing.T) {
	t.SkipNow()
}
