/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	IsAsyncCopy bool
	BatchSize   int
	BatchWait   int64
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

		if isDebug {
			bar.WithDisabled(false)
			bar64.WithDisabled64(false)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		Colorintln("green", "copy is running ...")
		log.Println("is-overwrite:", IsOverwrite)

		_, finputInfo, _ := NewBufReader(Input)
		if finputInfo.IsDir() == false {
			CopyFile(Input, Output)
		} else {
			if IsAsyncCopy == true {
				CopyDir2(Input, Output)
			} else {
				CopyDir(Input, Output)
			}

		}

	},
}

func init() {
	rootCmd.AddCommand(copyCmd)
	copyCmd.Flags().BoolVar(&IsOverwrite, "overwrite", false, "overwrite the existing file")
	copyCmd.Flags().BoolVar(&IsAsyncCopy, "async-copy", true, "copy file parallely")
	copyCmd.Flags().IntVar(&BatchSize, "batch", 200, "use with --copy-async=true")
	copyCmd.Flags().Int64Var(&BatchWait, "batch-wait", 5, "1 ~ 3600 seconds,use with --copy-async=true")

	rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagRequired("output")
}
