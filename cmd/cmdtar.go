package cmd

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	filesMap        map[string]string
	fileCompression int
	Excludes        string
	ExcludeFiles    []string
)

// tar2Cmd represents the tar command
var tarCmd = &cobra.Command{
	Use:   "tar",
	Short: "tar --input=/the/folder/you/want/to/tar --output=/the/file/where/you/want/to/save.tar",
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
		Colorintln("green", "tar is running ...")
		TarDir2(Input, Output)
	},
}

func init() {
	rootCmd.AddCommand(tarCmd)
	tarCmd.Flags().StringVar(&Excludes, "excludes", "", "a file to define excluded files, line by line")

	rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagRequired("output")
}
