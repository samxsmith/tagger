package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/samxsmith/tagger"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	t.SkipNow()
}

func TestListTag(t *testing.T) {
	t.SkipNow()
}

func TestRemoveTag(t *testing.T) {
	t.SkipNow()
}

func TestUpdateFileTags(t *testing.T) {
	t.SkipNow()
}

type fileToTag struct {
	filename string
	tags     []string
}

type testQuery struct {
	description   string
	queryString   string
	expectedFiles []string
}

func TestTagAndQuery(t *testing.T) {
	os.Remove("./test_data/files_to_tag/.tags")
	/**
		ADD TAGS
	**/
	testTagsToAdd := []fileToTag{
		{
			filename: "rose_img",
			tags:     []string{"subject:flower"},
		},
		{
			filename: "tulip_img",
			tags:     []string{"subject:flower"},
		},
		{
			filename: "wedding_picture",
			tags:     []string{"occassion:wedding", "person:me", "person:spouse"},
		},
		{
			filename: "sea_holiday_img",
			tags:     []string{"occassion:holiday", "subject:sea", "person:me"},
		},
		{
			filename: "sub_dir/nan_and_grandad_pic",
			tags:     []string{"person:nan", "person:grandad"},
		},
	}

	for _, tt := range testTagsToAdd {

		filepath := filepath.Join("test_data/files_to_tag", tt.filename)
		err := tagger.AddTag(filepath, tt.tags)
		assert.Nil(t, err)
	}

	/**
		QUERY TAGGED FILES
	**/

	queriesToTest := []testQuery{
		{
			description:   "query for a single tag, and get all matches",
			queryString:   "subject:flower",
			expectedFiles: []string{"rose_img", "tulip_img"},
		},
		{
			description:   "query for a person and get all their photos",
			queryString:   "person:me",
			expectedFiles: []string{"sea_holiday_img", "wedding_picture"},
		},
		{
			description:   "query by wild card and match anything with prefix",
			queryString:   "person:*",
			expectedFiles: []string{"sea_holiday_img", "sub_dir/nan_and_grandad_pic", "wedding_picture"},
		},
		{
			description:   "query with multiple to get files matching all tags",
			queryString:   "person:me person:spouse",
			expectedFiles: []string{"wedding_picture"},
		},
	}

	for i, tt := range queriesToTest {
		t.Run(fmt.Sprintf("case %d) %s", i, tt.description), func(t *testing.T) {
			files, err := tagger.QueryForFilenames(tt.queryString, nil)
			assert.Nil(t, err)
			for i, f := range tt.expectedFiles {
				tt.expectedFiles[i] = filepath.Join("test_data/files_to_tag", f)
			}
			sort.Strings(tt.expectedFiles)
			sort.Strings(files)
			assert.Equal(t, tt.expectedFiles, files)
		})
	}

}
