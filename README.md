# Tagger

An OS-independent file tagging mechanism to make file & photo storage useful, portable, structured & lightweight.

## Aims
Files & photos can be:
- tagged with multiple values
- queried by tags
- with a solution that can be backed up with files themselves
- in a lightweight and portable way
- that scales to large photo or file collections

## Install
`go install github.com/samxsmith/tagger`

## Usage
- `tagger init` to initialise the current directory with a `.tags` file
- `tagger add $FILENAME tag1 tag2 tag3 ...` to add tags to a file
- `tagger get-files tag1 tag3` to query for all files with that tag, recursively through the file structure from the current directory
- `tagger list-tags` to list all tags recursively.

## Example
```
photos/
  wedding_photo.jpg
  holiday_photo.jpg
  sub_dir/
    dog_and_me.jpg
    flower.jpg
```

Add tags
```shell
$ cd photos
$ tagger init
$ tagger add wedding_photo.jpg person:me person:spouse
$ tagger add holiday_photo person:spouse location:italy
$ tagger add sub_dir/dog_and_me.jpg location:home person:me person:dog
$ tagger add sub_dir/flower.jpg location:home
```

Query for files
```shell
$ tagger gf person:me

# photos with me in
wedding_photo.jpg
sub_dir/dog_and_me.jpg


$ tagger gf person:*

# photos with anyone in
wedding_photo.jpg
holiday_photo.jpg
sub_dir/dog_and_me.jpg

$ tagger gf person:spouse person:me

# photos with me & spouse in
wedding_photo.jpg
```
