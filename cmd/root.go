/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"time"

	"github.com/spf13/cobra"
)

var (
	Input    string
	Output   string
	BufferMB int
	tStart   time.Time
	tStop    time.Time
	isDebug  bool
)

var (
	numCPU int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gopgzip",
	Short: "-",
	Long:  `-`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		tStart = time.Now()

		Input = PathNormalize(Input)
		Output = PathNormalize(Output)

		if BufferMB < 0 || BufferMB > 2048 {
			BufferMB = 64
		}

		if Output != "" && strings.Index(Output, "://") < 0 {
			MakeDirs(filepath.Dir(Output))
		}

	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		tStop = time.Now()
		log.Printf("duration: %v sec", tStop.Sub(tStart))

		if isDebug {
			log.Println("-------------")
			log.Println("input:", Input)
			fi, err := os.Stat(Input)
			if err == nil && fi.IsDir() == false {
				log.Println("input xxhash:", Xxh3SumFile(Input))
			}

			log.Println("-------------")
			log.Println("output:", Output)
			fo, err := os.Stat(Output)
			if err == nil && fo.IsDir() == false {
				log.Println("output xxhash:", Xxh3SumFile(Output))
			}
		}
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
	numCPU = runtime.NumCPU()
	runtime.LockOSThread()
	runtime.GOMAXPROCS(numCPU)

	//
	rootCmd.PersistentFlags().StringVar(&Input, "input", "", "source file or folder(only [tar] need a folder here)")
	rootCmd.PersistentFlags().StringVar(&Output, "output", "", "target file or folder(only [untar] need a folder here)")
	rootCmd.PersistentFlags().IntVar(&BufferMB, "buffer-mb", 64, "1~2048;SSD: greater is better, HDD: lower is better")
	rootCmd.PersistentFlags().BoolVar(&isDebug, "debug", false, "will show more info if true")

}
