/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/harryzhu/pbar"

	"time"

	"github.com/spf13/cobra"
)

var (
	Input       string
	Output      string
	BufferMB    int
	tStart      time.Time
	tStop       time.Time
	IsOverwrite bool
	isDebug     bool
	isSIMD      bool
	isBar       bool
)

var (
	NumCPU int
	bar    *pbar.Bar
	bar64  *pbar.Bar64
)

const (
	MB int64 = int64(1024 * 1024)
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gopgzip",
	Short: "-",
	Long:  `-`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		NumCPU = runtime.NumCPU()
		runtime.LockOSThread()
		runtime.GOMAXPROCS(NumCPU)

		bar = pbar.NewBar(0)
		bar64 = pbar.NewBar64(0)

		Input = PathNormalize(Input)
		Output = PathNormalize(Output)

		if BufferMB < 0 || BufferMB > 2048 {
			BufferMB = 64
		}

		if Output != "" && strings.Index(Output, "://") < 0 {
			MakeDirs(filepath.Dir(Output))
		}

		tStart = time.Now()
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		tStop = time.Now()
		log.Printf("duration: %v sec", tStop.Sub(tStart))

		if isDebug {
			fmt.Println("*** Debug Info ***")
			log.Println("-------------")
			log.Println("input:", Input)

			log.Println("-------------")
			log.Println("output:", Output)

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

	//
	rootCmd.PersistentFlags().StringVar(&Input, "input", "", "source file or folder(only [tar] need a folder here)")
	rootCmd.PersistentFlags().StringVar(&Output, "output", "", "target file or folder(only [untar] need a folder here)")
	rootCmd.PersistentFlags().BoolVar(&IsOverwrite, "overwrite", false, "if overwrite the existing file")
	rootCmd.PersistentFlags().IntVar(&BufferMB, "buffer-mb", 64, "1~2048;SSD: greater is better, HDD: lower is better")
	rootCmd.PersistentFlags().BoolVar(&isDebug, "debug", false, "will show more info if true")
	rootCmd.PersistentFlags().BoolVar(&isSIMD, "simd", false, "use simd instructions")
	rootCmd.PersistentFlags().BoolVar(&isBar, "bar", false, "show progressbar")
}
