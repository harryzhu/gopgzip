package cmd

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	Excludes     string
	ExcludeFiles []string
)

// tar2Cmd represents the tar2 command
var tar2Cmd = &cobra.Command{
	Use:   "tar2",
	Short: "tar2 --input=/the/folder/you/want/to/tar --output=/the/file/where/you/want/to/save.tar",
	Long:  `--input= is a folder, --output= is a .tar file`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if Input == "" {
			log.Fatal("--input= cannot be empty and must be a folder")
		}

		Input = strings.TrimRight(Input, "/")
		Input = strings.TrimRight(Input, "\\")

		if Output == "" {
			Output = strings.Join([]string{filepath.Base(Input), "tar"}, ".")
			Output = filepath.Join(filepath.Dir(Input), Output)
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		Colorintln("green", "tar2 is running ...")
		TarDir2(Input, Output)
	},
}

func init() {
	rootCmd.AddCommand(tar2Cmd)
	tar2Cmd.Flags().StringVar(&Excludes, "excludes", "", "a file to define excluded files, line by line")

	rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagRequired("output")
}
