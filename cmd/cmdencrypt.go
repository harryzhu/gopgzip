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

var (
	Password string
	Force    bool
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "encrypt --input=original-file --output=encrypted-file [--password= --force]",
	Long:  passwordTips,
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

		PasswordTips()

		setKeyPasswordIV()
	},
	Run: func(cmd *cobra.Command, args []string) {
		Colorintln("green", "encrypt is running ...")

		AESEncodeFile(Input, Output)

	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)
	encryptCmd.Flags().StringVar(&Password, "password", "", "password for encrypt")
	encryptCmd.Flags().BoolVar(&Force, "force", false, "force")

	rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagRequired("output")
}
