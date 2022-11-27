/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	//"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "encrypt --input=original-file --output=encrypted-filename]",
	Long:  `-`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if Input == "" {
			log.Fatal("--input= cannot be empty")
		}

		if Output == "" {
			inputExt := filepath.Ext(Input)
			if inputExt == "" {
				Output = strings.Join([]string{Input, "enc"}, ".")
			} else {
				Output = strings.Replace(Input, inputExt, ".enc"+inputExt, 1)
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		Colorintln("green", "encrypt is running ...")

		AESEncodeFile(Input, Output)

	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)
	rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagRequired("output")
}
