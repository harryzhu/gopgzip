/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var (
	Input    string
	Output   string
	bufferMB int
	tStart   time.Time
	tStop    time.Time
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gopgzip",
	Short: "-",
	Long:  `-`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		tStart = time.Now()
		if bufferMB < 0 || bufferMB > 2048 {
			bufferMB = 8
		}

		if Output != "" {
			MakeDirs(filepath.Dir(Output))
		}

	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		tStop = time.Now()
		log.Printf("duration: %v sec", tStop.Sub(tStart))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&Input, "input", "", "input file you want to zip")
	rootCmd.PersistentFlags().StringVar(&Output, "output", "", "output file you want to save")
	rootCmd.PersistentFlags().IntVar(&bufferMB, "buffer-mb", 64, "1~2048,must less than memory available|SSD: greater is better, HDD: lower is better")
}
