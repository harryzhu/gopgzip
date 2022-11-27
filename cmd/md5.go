/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// md5Cmd represents the md5 command
var md5Cmd = &cobra.Command{
	Use:   "md5",
	Short: "md5 --input=your-local-file.txt [--output=if-you-want-to-save-hash-value-to-file.txt]",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		Colorintln("green", "md5 is running ...")
		m := MD5File(Input)
		Colorintln("green", m)

		if Output != "" {
			SaveFile(Output, []byte(m))
		}
	},
}

func init() {
	rootCmd.AddCommand(md5Cmd)
	rootCmd.MarkFlagRequired("input")

}
