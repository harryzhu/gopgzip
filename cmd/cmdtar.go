/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
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
	Long: `--input= is a folder, --output= is a file,
	--compression if you select 2(zstd) to tar, you need to select 2(zstd) to untar`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if Input == "" {
			log.Fatal("--input= cannot be empty")
		}

		if Output == "" {
			Output = strings.Join([]string{Input, "tar"}, ".")
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		Colorintln("green", "tar is running ...")
		TarballDir(Input, Output)

	},
}

func init() {
	rootCmd.AddCommand(tarCmd)
	tarCmd.Flags().IntVar(&fileCompression, "compression", 0, "0=None,1=gzip,2=zstd")
	rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagRequired("output")
}
