/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
)

var (
	IsOverwrite   bool
	IsKeepUrlPath bool
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "download --input=a-text-file-with-download-url-line-by-line.txt --output=download-root-dir",
	Long: ` --overwrite=true|false --keep-url-path=true|false;
	in "a-text-file-with-download-url-line-by-line.txt", 
	you can use "#" prefix to comment the line for ignoring the download,
	you can use only url in evey line, 
	or you can use "|" to split the soure_url and the local_save_path, i.e.:
	
	http://your-domain.com/static/file/1.zip
	http://your-domain.com/static/file/2.log
	http://your-domain.com/static/file/3.jpg | /home/harry/temp/333.jpg
	#http://your-domain.com/you-do-not-want-to-download/this-file.zip
	
	placeholder of the output,i.e.:
	--output=hello-world-{hostname}-{yyyy}{mm}{dd}-{HH}{MM}{SS}-{day-of-week}.png
	or:
	http://your-domain.com/static/file/3.jpg | /home/harry/temp/333-{hostname}-{yyyy}{mm}{dd}-{HH}{MM}{SS}-{day-of-week}.jpg
	`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if Input == "" || Output == "" {
			log.Fatal("--input= and --output= cannot be empty,")
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		Colorintln("green", "download is running ...")
		if strings.Index(Input, "http") == 0 {
			DownloadFile(Input, Output)
		} else {
			DownloadByList(Input, Output)
		}

	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().BoolVar(&IsOverwrite, "overwrite", false, "overwrite the existing file")
	downloadCmd.Flags().BoolVar(&IsKeepUrlPath, "keep-url-path", true, "keep-url-path")

	rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagRequired("output")
}
