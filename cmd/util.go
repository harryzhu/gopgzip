package cmd

import (
	"bufio"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	//"sync"
	//"time"

	//"github.com/klauspost/compress/zstd"
	gzip "github.com/klauspost/pgzip"
	"github.com/mholt/archiver/v4"
	progressbar "github.com/schollz/progressbar/v3"
	"github.com/zeebo/blake3"
)

func CompressZip(src, dst string) {
	numCPU := runtime.NumCPU()
	runtime.LockOSThread()
	runtime.GOMAXPROCS(numCPU)

	var selectNumCPU int = 1

	if numCPU > 1 && numCPU <= 4 {
		selectNumCPU = 2
	}

	if numCPU > 4 && numCPU <= 8 {
		selectNumCPU = 4
	}

	if numCPU > 8 {
		selectNumCPU = numCPU - 4
	}

	if Threads <= numCPU && Threads > 0 {
		selectNumCPU = Threads
	}

	var cLevel int = 6
	switch Level {
	case 0:
		cLevel = gzip.NoCompression
	case 1:
		cLevel = gzip.BestSpeed
	case 6:
		cLevel = gzip.DefaultCompression
	case 9:
		cLevel = gzip.BestCompression
	default:
		cLevel = gzip.DefaultCompression
	}

	fsrc, fsrcInfo, fsrcHandler := NewBufReader(src)

	dstTemp := strings.Join([]string{dst, "ing"}, "")

	fdst, fh := NewBufWriter(dstTemp)

	w, err := gzip.NewWriterLevel(fdst, cLevel)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	var BlockSizeByte int = BlockSizeMB << 20

	w.SetConcurrency(BlockSizeByte, selectNumCPU)

	log.Printf("threads: %v, block-size: %v MB", selectNumCPU, BlockSizeMB)

	if isDebug {
		bar := progressbar.DefaultBytes(fsrcInfo.Size())
		_, err = io.Copy(io.MultiWriter(w, bar), fsrc)
		bar.Finish()
	} else {
		_, err = io.Copy(w, fsrc)
	}
	if err != nil {
		log.Fatal(err)
	}

	w.Close()
	fdst.Flush()
	fh.Close()
	fsrcHandler.Close()

	_, err = os.Stat(dst)
	if err == nil {
		err = os.Remove(dst)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = os.Rename(dstTemp, dst)
	if err != nil {
		log.Fatal(err)
	}

	fullpathDst, _ := filepath.Abs(dst)

	Colorintln("green", "file: "+fullpathDst+"\n")
}

func DecompressZip(src string, dst string) error {
	fsrc, fsrcInfo, fhsrc := NewBufReader(src)

	dstTemp := strings.Join([]string{dst, "unzipping"}, ".")
	fdst, fhdst := NewBufWriter(dstTemp)

	reader, err := gzip.NewReader(fsrc)
	if err != nil {
		log.Fatal(err)
	}

	if isDebug {
		bar := progressbar.DefaultBytes(fsrcInfo.Size(), "unzipping ...")
		_, err = reader.WriteTo(io.MultiWriter(fdst, bar))
		bar.Finish()
	} else {
		_, err = reader.WriteTo(fdst)
	}

	if err != nil {
		log.Fatal(err)
	}
	fhdst.Close()
	fhsrc.Close()

	err = os.Rename(dstTemp, dst)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func MD5File(src string) string {
	reader, _, fhsrc := NewBufReader(src)
	hash := md5.New()

	var buf []byte = make([]byte, bufferMB)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		hash.Write(buf[:n])
	}

	fhsrc.Close()

	return hex.EncodeToString(hash.Sum(nil))
}

func SHA256File(src string) string {
	reader, _, fhsrc := NewBufReader(src)
	hash := sha256.New()

	var buf []byte = make([]byte, bufferMB)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		hash.Write(buf[:n])
	}

	fhsrc.Close()

	return hex.EncodeToString(hash.Sum(nil))
}

func SaveFile(src string, data []byte) error {
	err := ioutil.WriteFile(src, data, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func Blake3SumFile(src string) string {
	hash := blake3.New()
	reader, _, fhsrc := NewBufReader(src)

	var buf []byte = make([]byte, bufferMB)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		hash.Write(buf[:n])
	}

	fhsrc.Close()
	return hex.EncodeToString(hash.Sum(nil))
}

func Colorintln(c string, s string) error {
	s = strings.Join([]string{s, "\n"}, "")
	Colorint(c, s)
	return nil
}

func Colorint(c string, s string) error {
	platform := runtime.GOOS

	if platform == "windows" {
		fmt.Print(s)
		return nil
	}

	switch c {
	case "red":
		{
			fmt.Printf("\033[1;31;40m%s\033[0m\n", s)
		}
	case "green":
		{
			fmt.Printf("\033[1;32;40m%s\033[0m\n", s)
		}
	case "yellow":
		{
			fmt.Printf("\033[1;33;40m%s\033[0m\n", s)
		}
	default:
		{
			fmt.Printf("\033[1;31;40m%s\033[0m\n", s)
		}
	}
	return nil
}

func setFilesMap(src string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		log.Fatal(err)
	}
	if !srcInfo.IsDir() {
		log.Fatal(src + " should be a folder")
	}

	src = strings.ReplaceAll(src, "\\", "/")
	src = strings.TrimRight(src, "/")
	filesMap = make(map[string]string, 100)

	var walkFunc = func(path string, info os.FileInfo, err error) error {
		path = strings.ReplaceAll(path, "\\", "/")
		if !info.IsDir() {
			filesMap[path] = strings.Trim(strings.Replace(path, src[:strings.LastIndex(src, "/")], "", 1), "/")
		}

		return nil
	}
	err = filepath.Walk(src, walkFunc)
	return err
}

func TarballDir(src string, dst string) error {
	setFilesMap(src)
	files, err := archiver.FilesFromDisk(nil, filesMap)
	if err != nil {
		log.Fatal(err)
	}

	dstTemp := strings.Join([]string{dst, "ing"}, "")

	bufdst, fhdst := NewBufWriter(dstTemp)

	defer func() {
		bufdst.Flush()
	}()

	format := archiver.CompressedArchive{
		Compression: nil,
		Archival:    archiver.Tar{},
	}
	if fileCompression == 1 {
		format = archiver.CompressedArchive{
			Compression: archiver.Gz{},
			Archival:    archiver.Tar{},
		}
	}

	if fileCompression == 2 {
		format = archiver.CompressedArchive{
			Compression: archiver.Zstd{},
			Archival:    archiver.Tar{},
		}
	}

	if isDebug {
		bar := progressbar.DefaultBytes(-1)
		err = format.Archive(context.Background(), io.MultiWriter(bufdst, bar), files)
		if err != nil {
			log.Fatal(err)
		}

		bufdst.Flush()
		fhdst.Close()
		bar.Finish()
	} else {
		err = format.Archive(context.Background(), bufdst, files)
		if err != nil {
			log.Fatal(err)
		}

		bufdst.Flush()
		fhdst.Close()
	}

	err = os.Rename(dstTemp, dst)
	if err != nil {
		log.Fatal(err)
	}

	fullpathDst, _ := filepath.Abs(dst)
	Colorintln("green", "file: "+fullpathDst)

	return nil
}

func Untarball(src string, dst string) error {
	fsrc, _, fhsrc := NewBufReader(src)
	wg := sync.WaitGroup{}

	format := archiver.CompressedArchive{
		Compression: nil,
		Archival:    archiver.Tar{},
	}

	if fileCompression == 1 {
		format = archiver.CompressedArchive{
			Compression: archiver.Gz{},
			Archival:    archiver.Tar{},
		}
	}

	if fileCompression == 2 {
		format = archiver.CompressedArchive{
			Compression: archiver.Zstd{},
			Archival:    archiver.Tar{},
		}
	}

	handler := func(ctx context.Context, f archiver.File) error {
		rc, err := f.Open()
		if err != nil {
			log.Println(err)
			return err
		}
		defer rc.Close()

		srcStat, err := f.Stat()
		if err != nil {
			log.Println(err)
		}
		srcData, err := io.ReadAll(rc)
		dstName := filepath.Join(dst, f.NameInArchive)
		MakeDirs(filepath.Dir(dstName))

		//fmt.Println(dstName)
		if srcStat.Size() > 32<<20 {
			err = ioutil.WriteFile(dstName, srcData, srcStat.Mode())
			if err != nil {
				log.Println(err)
			}
		} else {
			wg.Add(1)
			go func() {
				err = ioutil.WriteFile(dstName, srcData, srcStat.Mode())
				if err != nil {
					log.Println(err)
				}
				wg.Done()
			}()
		}
		return err
	}

	err := format.Extract(context.Background(), fsrc, nil, handler)
	if err != nil {
		log.Fatal(err)
	}
	fhsrc.Close()

	wg.Wait()
	return nil
}

func MakeDirs(s string) error {
	_, err := os.Stat(s)
	if err == os.ErrExist {
		return nil
	}
	err = os.MkdirAll(s, os.ModePerm)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func NewBufWriter(f string) (*bufio.Writer, *os.File) {
	fh, err := os.Create(f)
	if err != nil {
		log.Fatal(err)
	}

	return bufio.NewWriterSize(bufio.NewWriter(fh), bufferMB), fh
}

func NewBufReader(f string) (*bufio.Reader, fs.FileInfo, *os.File) {
	finfo, err := os.Stat(f)
	if err != nil {
		log.Fatal(err)
	}

	fh, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}

	return bufio.NewReaderSize(bufio.NewReader(fh), bufferMB), finfo, fh
}
