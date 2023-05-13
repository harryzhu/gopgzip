/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// md5Cmd represents the md5 command
var md5Cmd = &cobra.Command{
	Use:   "md5",
	Short: "md5 --input=your-local-file.txt [--output=if-you-want-to-save-hash-value-to-file.txt]",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		m := ""
		if isSIMD {
			log.Println("Support SIMD 2")
			m = MD5FileSIMD(Input)
		} else {
			m = MD5File(Input)
		}

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
