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
	filesMap map[string]string
)

// tarCmd represents the tar command
var tarCmd = &cobra.Command{
	Use:   "tar",
	Short: "tar --input=/the/root/dir/you/want/to/tar --output=/the/path/where/you/want/to/save.tar",
	Long:  `--input= is a folder, --output= is a file.`,
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

	rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagRequired("output")
}
