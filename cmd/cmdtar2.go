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
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// tar2Cmd represents the tar2 command
var tar2Cmd = &cobra.Command{
	Use:   "tar2",
	Short: "tar2 --input=/the/folder/you/want/to/tar --output=/the/file/where/you/want/to/save.tar",
	Long:  `--input= is a folder, --output= is a .tar file`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if Input == "" {
			log.Fatal("--input= cannot be empty and must be a folder")
		}

		Input = strings.TrimRight(Input, "/")
		Input = strings.TrimRight(Input, "\\")

		if Output == "" {
			Output = strings.Join([]string{filepath.Base(Input), "tar"}, ".")
			Output = filepath.Join(filepath.Dir(Input), Output)
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		Colorintln("green", "tar2 is running ...")
		TarDir2(Input, Output)
	},
}

func init() {
	rootCmd.AddCommand(tar2Cmd)

	rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagRequired("output")
}
