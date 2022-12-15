package cmd

import (
	"archive/tar"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/harryzhu/pbar"
)

func TarDir2(src string, dst string) error {
	src, _ = filepath.Abs(src)
	src = strings.TrimRight(src, "/")

	setFilesMap(src)
	dstTemp := strings.Join([]string{dst, "ing"}, "")
	bufdst, fhdst := NewBufWriter(dstTemp)

	tw := tar.NewWriter(bufdst)
	bar := pbar.NewBar64(0)

	var bufSize int64 = 0
	var bufByte int64 = int64(BufferMB << 20)
	var fsrc io.Reader
	var fsrcInfo fs.FileInfo
	var fhsrc *os.File
	var hdr *tar.Header
	var err error

	for fk, fv := range filesMap {
		fsrc, fsrcInfo, fhsrc = NewBufReader(fk)

		hdr, err = tar.FileInfoHeader(fsrcInfo, fv)
		if err != nil {
			log.Fatal(err)
		}

		hdr.Name = strings.TrimPrefix(fv, "/")

		err = tw.WriteHeader(hdr)
		if err != nil {
			log.Fatal(err)
		}

		if isDebug {
			_, err = io.Copy(io.MultiWriter(tw, bar), fsrc)
		} else {
			_, err = io.Copy(tw, fsrc)
		}

		if err != nil {
			log.Fatal(err)
		}

		fhsrc.Close()

		bufSize += fsrcInfo.Size()
		if bufByte-bufSize < 1024 {
			if isDebug {
				log.Println("Flush: ", bufSize)
			}
			bufdst.Flush()
			bufSize = 0
		}
	}

	tw.Close()
	bufdst.Flush()
	fhdst.Close()

	err = os.Rename(dstTemp, dst)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func UntarDir2(src string, dst string) error {
	dst, _ = filepath.Abs(dst)
	dst = strings.TrimRight(dst, "/")

	fsrc, _, fhsrc := NewBufReader(src)

	tr := tar.NewReader(fsrc)
	bar := pbar.NewBar64(0)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fname := filepath.Join(dst, hdr.Name)

		if hdr.FileInfo().IsDir() {
			MakeDirs(fname)
			continue
		}

		MakeDirs(filepath.Dir(fname))

		fdst, fhdst := NewBufWriter(fname)
		if isDebug {
			_, err = io.Copy(io.MultiWriter(fdst, bar), tr)
		} else {
			_, err = io.Copy(fdst, tr)
		}

		if err != nil {
			log.Fatal(err)
		}
		fhdst.Close()
		//
		err = os.Chtimes(fname, hdr.AccessTime, hdr.ModTime)
		if err != nil {
			log.Println(err)
		}

		err = os.Chmod(fname, hdr.FileInfo().Mode())
		if err != nil {
			log.Println(err)
		}

		err = os.Chown(fname, hdr.Uid, hdr.Gid)
		if err != nil {
			log.Println(err)
		}
	}

	bar.Finish()

	fhsrc.Close()

	return nil
}
