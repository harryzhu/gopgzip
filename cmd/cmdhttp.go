/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

var (
	HttpIP   string
	HttpPort string
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "http --ip=your-machine-ip --port=8080 --input=the-absolute-path-you-want-to-serve",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if Input == "" {
			Input = "./"
		}

		HttpServer(HttpIP, HttpPort, Input)
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
	httpCmd.Flags().StringVar(&HttpIP, "ip", "0.0.0.0", "the machine's ip address")
	httpCmd.Flags().StringVar(&HttpPort, "port", "8080", "the port, default port is 8080")

	rootCmd.MarkFlagRequired("input")
}
