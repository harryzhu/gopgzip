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

// unzstdCmd represents the unzstd command
var unzstdCmd = &cobra.Command{
	Use:   "unzstd",
	Short: "unzstd --input=your-backup.zst --output=your-local-file.txt",
	Long:  `-`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if Input == "" {
			log.Fatal("--input= cannot be empty")
		}

		if Output == "" {
			log.Println("you can use --output=your-filepath for saving")
			Output = strings.Replace(Input, ".zst", "", 1)
			_, err := os.Stat(Output)
			if err == nil {
				t := time.Now().Format("20060102150405")
				Output = strings.Replace(Output, filepath.Ext(Output), "_"+t+filepath.Ext(Output), -1)
			}
			log.Println("the default file(in the same folder): " + Output)
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		Colorintln("green", "unzstd is running ...")

		DecompressWithZstd(Input, Output)
	},
}

func init() {
	rootCmd.AddCommand(unzstdCmd)
}
