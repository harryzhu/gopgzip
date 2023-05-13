package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// sha256Cmd represents the sha256 command
var sha256Cmd = &cobra.Command{
	Use:   "sha256",
	Short: "sha256 --input=your-local-file.txt [--output=if-you-want-to-save-hash-value-to-file.txt]",
	Long:  `-`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if Input == "" {
			log.Fatal("--input= cannot be empty")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		p := ""

		if isAvx512 {
			log.Println("Support AVX512")
			p = SHA256FileSIMD(Input)
		} else {
			p = SHA256File(Input)
		}

		Colorintln("green", p)

		if Output != "" {
			SaveFile(Output, []byte(p))
		}
	},
}

func init() {
	rootCmd.AddCommand(sha256Cmd)

	rootCmd.MarkFlagRequired("input")
}
