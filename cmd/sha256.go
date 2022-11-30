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
	"github.com/spf13/cobra"
)

// sha256Cmd represents the sha256 command
var sha256Cmd = &cobra.Command{
	Use:   "sha256",
	Short: "sha256 --input=your-local-file.txt [--output=if-you-want-to-save-hash-value-to-file.txt]",
	Long:  `-`,
	Run: func(cmd *cobra.Command, args []string) {
		Colorintln("green", "sha256 is running ...")
		m := SHA256File(Input)
		Colorintln("green", m)

		if Output != "" {
			SaveFile(Output, []byte(m))
		}
	},
}

func init() {
	rootCmd.AddCommand(sha256Cmd)

	rootCmd.MarkFlagRequired("input")
}
