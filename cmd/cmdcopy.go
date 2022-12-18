/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "A brief description of your command",
	Long:  `-`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("copy called")

		//CopyFile(Input, Output)
		CopyDir(Input, Output)
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)

}
