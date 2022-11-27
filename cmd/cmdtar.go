/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	//"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	filesMap map[string]string
	bufferMB int
)

// tarCmd represents the tar command
var tarCmd = &cobra.Command{
	Use:   "tar",
	Short: "tar --input=your-dir --output=/the/path/where/you/want/to/save.tar",
	Long:  ``,
	PreRun: func(cmd *cobra.Command, args []string) {
		if Input == "" || Output == "" {
			log.Fatal("--input= and --output= cannot be empty")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		Colorintln("green", "tar is running ...")
		TarballDir(Input, Output)

	},
}

func init() {
	rootCmd.AddCommand(tarCmd)
	tarCmd.Flags().IntVar(&bufferMB, "buffer-mb", 8, "buffer size: 1~1024")

	rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagRequired("output")
}
