package cmd

import (
	"archive/tar"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"time"
)

func TarDir2(src string, dst string) error {
	src, _ = filepath.Abs(src)
	src = strings.TrimRight(src, "/")

	setFilesMap(src)
	dstTemp := strings.Join([]string{dst, "ing"}, "")
	bufdst, fhdst := NewBufWriter(dstTemp)

	tw := tar.NewWriter(bufdst)
	bar64.WithMax64(0)

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
		hdr.ModTime = fsrcInfo.ModTime()

		err = tw.WriteHeader(hdr)
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(io.MultiWriter(tw, bar64), fsrc)

		if err != nil {
			log.Fatal(err)
		}

		fhsrc.Close()

		bufSize += fsrcInfo.Size()
		if bufByte-bufSize < 1024 {
			bufdst.Flush()
			bufSize = 0
		}
	}
	tw.Flush()
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
	bar64.WithMax64(0)

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

		_, err = io.Copy(io.MultiWriter(fdst, bar64), tr)

		if err != nil {
			log.Fatal(err)
		}
		fdst.Flush()
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

	bar64.Finish()

	fhsrc.Close()

	return nil
}

func CopyFile2(src string, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		log.Fatal(err)
	}

	dstFile.Close()
	srcFile.Close()
	return nil
}

func CopyDir2(src string, dst string) error {
	if BatchWait < 1 || BatchWait > 3600 {
		BatchWait = 3
	}
	wgc := sync.WaitGroup{}
	bar.WithDisabled(false)
	bar.WithMax(0)
	var copyPool int

	var walkFunc = func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		srcPath, _ := filepath.Abs(path)
		dstPath := strings.Replace(srcPath, filepath.Dir(src), dst, 1)

		MakeDirs(filepath.Dir(dstPath))

		if IsOverwrite == false {
			_, err := FileInfo(dstPath)
			if err == nil {
				return nil
			}
		}

		if copyPool >= BatchSize {
			time.Sleep(time.Second * time.Duration(BatchWait))
		}

		if isDebug {
			log.Println("copyPool:", copyPool)
		}

		wgc.Add(1)

		go func() {
			bar.Add(1)
			CopyFile2(srcPath, dstPath)
			if copyPool > 0 {
				copyPool -= 1
			}
			wgc.Done()
		}()

		copyPool += 1

		return err
	}

	filepath.Walk(src, walkFunc)

	wgc.Wait()

	return nil
}
