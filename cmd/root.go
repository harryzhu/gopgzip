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

	"github.com/harryzhu/pbar"

	"time"

	. "github.com/klauspost/cpuid/v2"
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
	isAvx512    bool
)

var (
	NumCPU int
	bar    *pbar.Bar
	bar64  *pbar.Bar64
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gopgzip",
	Short: "-",
	Long:  `-`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		isSIMD = false
		isAvx512 = false

		if CPU.Supports(SSE, SSE2) {
			isSIMD = true
		}
		//if CPU.Supports(SHA, SSSE3, SSE4) {
		if CPU.Supports(AVX512F, AVX512DQ, AVX512BW, AVX512VL) {
			isAvx512 = true
		}

		NumCPU = runtime.NumCPU()
		runtime.LockOSThread()
		runtime.GOMAXPROCS(NumCPU)

		bar = pbar.NewBar(0)
		bar64 = pbar.NewBar64(0)

		bar.WithDisabled(true)
		bar64.WithDisabled64(true)

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

	//
	rootCmd.PersistentFlags().StringVar(&Input, "input", "", "source file or folder(only [tar] need a folder here)")
	rootCmd.PersistentFlags().StringVar(&Output, "output", "", "target file or folder(only [untar] need a folder here)")
	rootCmd.PersistentFlags().BoolVar(&IsOverwrite, "overwrite", false, "if overwrite the existing file")
	rootCmd.PersistentFlags().IntVar(&BufferMB, "buffer-mb", 64, "1~2048;SSD: greater is better, HDD: lower is better")
	rootCmd.PersistentFlags().BoolVar(&isDebug, "debug", false, "will show more info if true")

}
