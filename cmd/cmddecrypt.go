/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "decrypt --input=your-encrypted-file --output=/where/you/want/to/decrypt/filename",
	Long:  `-`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if Input == "" {
			log.Fatal("--input= cannot be empty")
		}

		if Output == "" {
			inputExt := filepath.Ext(Input)
			if inputExt == "" {
				Output = strings.Join([]string{Input, "dec"}, ".")
			} else {
				Output = strings.Replace(Input, inputExt, ".dec"+inputExt, 1)
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		Colorintln("green", "decrypt is running ...")
		AESDecodeFile(Input, Output)
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)
	rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagRequired("output")
}
