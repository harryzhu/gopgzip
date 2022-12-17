/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// xxh3sumCmd represents the xxh3sum command
var xxhashCmd = &cobra.Command{
	Use:   "xxhash",
	Short: "xxhash --input=abc.iso [--output=sum.txt]",
	Long:  ``,
	PreRun: func(cmd *cobra.Command, args []string) {
		if Input == "" {
			log.Fatal("--input= cannot be empty")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		p := Xxh3SumFile(Input)

		if Output != "" {
			SaveFile(Output, []byte(p))
		}

		Colorintln("green", p)
	},
}

func init() {
	rootCmd.AddCommand(xxhashCmd)

}
