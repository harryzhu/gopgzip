/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// zstdCmd represents the zstd command
var zstdCmd = &cobra.Command{
	Use:   "zstd",
	Short: "A brief description of your command",
	Long:  `-`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if Input == "" {
			log.Fatal("--input= cannot be empty")
		}

		if Output == "" {
			Output = Input + ".zst"
		}

		if BlockSizeMB < 0 || BlockSizeMB > 512 {
			BlockSizeMB = 16
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		Colorintln("green", "zstd is running ...")

		CompressWithZstd(Input, Output)
	},
}

func init() {
	rootCmd.AddCommand(zstdCmd)
	zstdCmd.Flags().IntVar(&Level, "level", 1, "level should be 0,1,6,9; default 6")

}