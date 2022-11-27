/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// b3sumCmd represents the b3sum command
var b3sumCmd = &cobra.Command{
	Use:   "b3sum",
	Short: "b3sum --input=your-local-file.txt [--output=if-you-want-to-save-hash-value-to-file.txt]",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		p := Blake3SumFile(Input)
		Colorintln("green", p)

		if Output != "" {
			SaveFile(Output, []byte(p))
		}
	},
}

func init() {
	rootCmd.AddCommand(b3sumCmd)
}
