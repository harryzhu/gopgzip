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

var passwordDefault string = "This(*Key*)@2021That[#Key$]&1202"

var passwordTips string = `
for security issue, do NOT use --password in command line. use env variables instead.
in your /etc/profile, add: export HARRYZHUENCRYPTKEY=Your-Password;
then open a new terminal window, run encypt or decrypt.
if you still want to use --password= in command line, use --force=true meanwhile.
ie.: 
encrypt --input=doc.tar.gz --output=doc.tar.gz.enc --password="123" --force
decrypt --input=doc.tar.gz.enc --output=doc.tar.gz --password="123" --force

if not any password set, will use a defalut password: 
` + passwordDefault

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "encrypt --input=original-file --output=encrypted-filename [--password= --force]",
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

		if Password != "" {
			Colorintln("yellow", passwordTips)
		}

		if Password != "" && Force == false {
			log.Fatal("--password= must be with --force")
		}

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
