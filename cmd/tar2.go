package cmd

import (
	"archive/tar"
	"io"
	"log"
	"path/filepath"

	"os"
	"strings"
)

func TarDir2(src string, dst string) error {
	src, _ = filepath.Abs(src)
	src = strings.TrimRight(src, "/")

	setFilesMap(src)
	dstTemp := strings.Join([]string{dst, "ing"}, "")
	bufdst, fhdst := NewBufWriter(dstTemp)

	tw := tar.NewWriter(bufdst)

	for fk, fv := range filesMap {
		//log.Println(fk)
		//log.Println(fv)

		fsrc, fsrcInfo, fhsrc := NewBufReader(fk)

		hdr, err := tar.FileInfoHeader(fsrcInfo, fv)
		if err != nil {
			log.Fatal(err)
		}
		hdr.Name = strings.TrimPrefix(fv, "/")

		err = tw.WriteHeader(hdr)
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(tw, fsrc)
		if err != nil {
			log.Fatal(err)
		}

		fhsrc.Close()

	}

	tw.Close()
	bufdst.Flush()
	fhdst.Close()

	err := os.Rename(dstTemp, dst)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
