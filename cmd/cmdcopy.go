/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "copy --input=a-folder-or-file  --output=a-folder-or-file",
	Long:  `-`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if Input == "" || Output == "" {
			log.Fatal("--input= and --output= cannot be empty")
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		Colorintln("green", "copy is running ...")

		_, finputInfo, _ := NewBufReader(Input)
		if finputInfo.IsDir() == false {
			CopyFile(Input, Output)
		} else {
			CopyDir(Input, Output)
		}

	},
}

func init() {
	rootCmd.AddCommand(copyCmd)

	rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagRequired("output")
}
