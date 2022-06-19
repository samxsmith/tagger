package main

import (
	"fmt"

	"github.com/samxsmith/tagger"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{}

var addTagsCmd = &cobra.Command{
	Use:     "add -f path/to/file.jpg tag1 [tag2 ....]",
	Short:   "Add tags to a file",
	Aliases: []string{"a"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Pass at least one tag")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		file, err := cmd.Flags().GetString("file")
		if err != nil {
			fmt.Println("ERROR: ", err)
			return
		}
		if err := tagger.AddTag(file, args); err != nil {
			fmt.Println("ERROR: ", err)
			return
		}
		fmt.Printf("ADDED TAGS <%v> to file `%s`", args, file)
	},
}

var getFilesCmd = &cobra.Command{
	Use:     "get-files tag1 [tag2 ...]",
	Short:   "Get the files matching all provided tags",
	Aliases: []string{"gf"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		files, err := tagger.QueryForFilenames(".", args, nil)
		if err != nil {
			fmt.Println("ERROR: ", err)
			return
		}
		for _, f := range files {
			fmt.Println(f)
		}
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise a .tags file in the current directory",
	Run: func(cmd *cobra.Command, args []string) {
		if err := tagger.Init("."); err != nil {
			fmt.Println("ERROR: ", err)
			return
		}
		fmt.Println(".tags file initialised")
	},
}

var listTagsCmd = &cobra.Command{
	Use:     "list-tags [prefix]",
	Short:   "List all tags recursively from current directory. Use a prefix arg to filter.",
	Aliases: []string{"lt"},
	Run: func(cmd *cobra.Command, args []string) {
		var prefix string
		if len(args) > 0 {
			prefix = args[0]
		}
		ls, err := tagger.ListTags(".", prefix)
		if err != nil {
			fmt.Println("ERROR: ", err)
			return
		}
		for _, f := range ls {
			fmt.Println(f)
		}
	},
}

func main() {
	// FLAGS
	addTagsCmd.Flags().StringP("file", "f", "", "The path to the file you want to tag")
	addTagsCmd.MarkFlagRequired("file")

	// CMDs
	rootCmd.AddCommand(addTagsCmd)
	rootCmd.AddCommand(getFilesCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(listTagsCmd)

	// EXEC
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("ERROR: ", err)
	}
}
