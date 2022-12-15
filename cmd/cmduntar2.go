/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// untar2Cmd represents the untar2 command
var untar2Cmd = &cobra.Command{
	Use:   "untar2",
	Short: "untar2 --input=your-file.tar --output=/the/folder/where/you/want/to/extract_dir",
	Long:  `--input= is a .tar file, --output= is a folder`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if Input == "" || Output == "" {
			log.Fatal("--input= and --output= cannot be empty")
		}

		_, err := os.Stat(Input)
		if err != nil {
			log.Fatal(err)
		}

		err = MakeDirs(Output)
		if err != nil {
			log.Fatal(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		Colorintln("green", "untar2 is running ...")
		UntarDir2(Input, Output)
	},
}

func init() {
	rootCmd.AddCommand(untar2Cmd)

	rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagRequired("output")
}
