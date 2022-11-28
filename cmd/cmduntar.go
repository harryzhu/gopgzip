/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	//"fmt"
	"log"

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

		_, err := os.Stat(Input)
		if err != nil {
			log.Fatal(err)
		}

		err = MakeDirs(Output)
		if err != nil {
			log.Fatal(err)
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
