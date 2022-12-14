/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// unzipCmd represents the unzip command
var unzipCmd = &cobra.Command{
	Use:   "unzip",
	Short: "unzip --input=your-local-gzip-file.gz [--output=unzip-to-local-disk-file]",
	Long:  `-`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if Input == "" {
			log.Fatal("--input= cannot be empty")
		}

		if Output == "" {
			log.Println("you can use --output=your-filepath for saving")
			Output = strings.Replace(Input, ".gz", "", 1)
			_, err := os.Stat(Output)
			if err == nil {
				t := time.Now().Format("20060102150405")
				Output = strings.Replace(Output, filepath.Ext(Output), "_"+t+filepath.Ext(Output), -1)
			}
			log.Println("the default file(in the same folder): " + Output)
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		Colorintln("green", "unzip is running ...")

		DecompressWithGZip(Input, Output)
	},
}

func init() {
	rootCmd.AddCommand(unzipCmd)
	rootCmd.MarkFlagRequired("input")
}
