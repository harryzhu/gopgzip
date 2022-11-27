/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	//"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// untarCmd represents the untar command
var untarCmd = &cobra.Command{
	Use:   "untar",
	Short: "untar --input=your-file.tar --output=/the/path/where/you/want/to/extract-folder",
	Long:  `-`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if Input == "" || Output == "" {
			log.Fatal("--input= and --output= cannot be empty")
		}
		finfo, err := os.Stat(Output)

		if err == nil && finfo.IsDir() == false {
			log.Fatal("--output= should be a folder, not a single file")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		Colorintln("green", "untar is running ...")
		Untarball(Input, Output)
	},
}

func init() {
	rootCmd.AddCommand(untarCmd)
	rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagRequired("output")

}
