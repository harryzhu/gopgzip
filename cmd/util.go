package cmd

import (
	"bufio"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	//"github.com/klauspost/compress/zstd"
	gzip "github.com/klauspost/pgzip"
	"github.com/mholt/archiver/v4"
	progressbar "github.com/schollz/progressbar/v3"
	"github.com/zeebo/blake3"
)

func CompressZip(src, dst string) {
	numCPU := runtime.NumCPU()

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

	runtime.LockOSThread()
	runtime.GOMAXPROCS(numCPU)

	tStart := time.Now()

	fsrc, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer fsrc.Close()

	dstTemp := strings.Join([]string{dst, "ing"}, "")
	fdst, err := os.Create(dstTemp)
	if err != nil {
		log.Fatal(err)
	}
	defer fdst.Close()

	finfoSrc, err := os.Stat(src)
	if err != nil {
		log.Fatal(err)
	}

	w, err := gzip.NewWriterLevel(fdst, cLevel)
	if err != nil {
		log.Fatal(err)
	}

	defer w.Close()

	var BlockSizeByte int = 16 << 20
	if BlockSizeMB > 0 {
		BlockSizeByte = BlockSizeMB << 20
	}

	w.SetConcurrency(BlockSizeByte, selectNumCPU)

	log.Printf("threads: %v, block-size: %v MB", selectNumCPU, BlockSizeMB)

	bar := progressbar.DefaultBytes(finfoSrc.Size())
	_, err = io.Copy(io.MultiWriter(w, bar), fsrc)
	if err != nil {
		log.Fatal(err)
	}
	bar.Finish()

	w.Close()
	fdst.Close()

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
	tStop := time.Now()
	duration := tStop.Sub(tStart)

	fullpathDst, _ := filepath.Abs(dst)

	log.Printf("OK. duration: %v sec\n", duration)
	Colorintln("green", "file: "+fullpathDst+"\n")
}

func DecompressZip(src string, dst string) error {
	fsrc, err := os.Open(src)
	defer fsrc.Close()
	if err != nil {
		log.Fatal(err)
	}
	dstTemp := strings.Join([]string{dst, "unzipping"}, ".")
	fdst, err := os.Create(dstTemp)
	defer fdst.Close()
	if err != nil {
		log.Fatal(err)
	}

	reader, err := gzip.NewReader(fsrc)
	if err != nil {
		log.Fatal(err)
	}

	bar := progressbar.DefaultBytes(-1, "unzipping ...")

	_, err = reader.WriteTo(io.MultiWriter(fdst, bar))
	if err != nil {
		log.Fatal(err)
	}

	err = os.Rename(dstTemp, dst)
	if err != nil {
		log.Fatal(err)
	}

	bar.Finish()
	return nil
}

func MD5File(src string) string {
	f, err := os.Open(src)
	defer f.Close()

	if err != nil {
		log.Fatal(err)
	}
	tStart := time.Now()

	hash := md5.New()

	var bufferSize = 64 << 20
	var buf []byte = make([]byte, bufferSize)

	reader := bufio.NewReader(f)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		//log.Println(n)
		hash.Write(buf[:n])

	}
	tStop := time.Now()
	log.Printf("duration: %v sec", tStop.Sub(tStart))
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
	f, err := os.Open(src)
	defer f.Close()

	if err != nil {
		log.Fatal(err)
	}
	tStart := time.Now()

	hash := blake3.New()

	finfo, err := os.Stat(src)
	if err != nil {
		log.Fatal(err)
	}
	fsize := finfo.Size()
	bar := progressbar.DefaultBytes(fsize)

	var bufferSize = 32 << 20
	var buf []byte = make([]byte, bufferSize)

	reader := bufio.NewReader(f)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		hash.Write(buf[:n])
		bar.Add64(int64(n))
	}
	bar.Finish()
	tStop := time.Now()

	log.Printf("duration: %v sec", tStop.Sub(tStart))
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
	default:
		{
			fmt.Printf("\033[1;31;40m%s\033[0m\n", s)
		}
	}
	return nil
}

func setFilesMap(src string) error {
	filesMap = make(map[string]string, 100)

	var walkFunc = func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			filesMap[path] = strings.Trim(strings.Replace(path, src, "", 1), "/")
		}

		return nil
	}
	err := filepath.Walk(src, walkFunc)
	return err
}

func TarballDir(src string, dst string) error {
	setFilesMap(src)
	files, err := archiver.FilesFromDisk(nil, filesMap)
	if err != nil {
		log.Fatal(err)
	}

	dstTemp := strings.Join([]string{dst, "ing"}, "")
	fdst, err := os.Create(dstTemp)
	if err != nil {
		log.Fatal(err)
	}
	defer fdst.Close()

	format := archiver.CompressedArchive{
		Compression: nil,
		Archival:    archiver.Tar{},
	}

	bar := progressbar.DefaultBytes(-1)

	err = format.Archive(context.Background(), io.MultiWriter(fdst, bar), files)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Rename(dstTemp, dst)
	if err != nil {
		log.Fatal(err)
	}

	bar.Finish()
	fullpathDst, _ := filepath.Abs(dst)
	Colorintln("green", "file: "+fullpathDst)

	return nil
}

func Untarball(src string, dst string) error {
	fsrc, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	bar := progressbar.DefaultBytes(-1)

	format := archiver.Tar{}

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
		dstName := filepath.Join(Output, f.NameInArchive)
		//log.Println(srcStat.Mode(), " ", dstName)

		MakeDirs(filepath.Dir(dstName))

		err = ioutil.WriteFile(dstName, srcData, srcStat.Mode())

		if err != nil {
			log.Println(err)
		}

		bar.Add64(srcStat.Size())
		return err
	}

	err = format.Extract(context.Background(), fsrc, nil, handler)
	if err != nil {
		log.Fatal(err)
	}
	bar.Finish()

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
	}
	return err
}
