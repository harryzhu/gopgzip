/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	BlockSizeMB int
	Level       int
	Threads     int
)

// zipCmd represents the zip command
var zipCmd = &cobra.Command{
	Use:   "zip",
	Short: "zip --input=your-local-file.txt --output=your-backup.gz [--level=0|1|6|9 --threads=8]",
	Long:  `-`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if Input == "" {
			log.Fatal("--input= cannot be empty")
		}

		if Output == "" {
			Output = Input + ".gz"
		}

		if BlockSizeMB < 0 || BlockSizeMB > 512 {
			BlockSizeMB = 16
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		Colorintln("green", "zip is running ...")

		CompressZip(Input, Output)
	},
}

func init() {
	rootCmd.AddCommand(zipCmd)
	zipCmd.Flags().IntVar(&BlockSizeMB, "block-size-mb", 16, "block size megabytes")
	zipCmd.Flags().IntVar(&Level, "level", 6, "level should be 0,1,6,9; default 6")
	zipCmd.Flags().IntVar(&Threads, "threads", 0, "threads for zip; default 0: auto-detect")

	rootCmd.MarkFlagRequired("input")
}
