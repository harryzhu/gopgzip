/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"io"
	"log"
)

func CopyFileWithBar(r io.Reader, w io.Writer, barMax64 int64) (err error) {
	if isBar {
		bar64.WithMax64(barMax64)
		_, err = io.Copy(io.MultiWriter(w, bar64), r)
		bar64.Finish()
	} else {
		_, err = io.Copy(w, r)
	}

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
