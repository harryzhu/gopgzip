/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"path/filepath"
	"strings"

	//"os"
	//"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	filesMap        map[string]string
	fileCompression int
)

// tarCmd represents the tar command
var tarCmd = &cobra.Command{
	Use:   "tar",
	Short: "tar --input=/the/folder/you/want/to/tar --output=/the/file/where/you/want/to/save.tar",
	Long:  `--input= is a folder, --output= is a file`,
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
		TarDir(Input, Output)

	},
}

func init() {
	rootCmd.AddCommand(tarCmd)
	rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagRequired("output")
}
