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

	"github.com/harryzhu/pbar"
	"github.com/klauspost/compress/zstd"
	gzip "github.com/klauspost/pgzip"
	"github.com/mholt/archiver/v4"

	//"github.com/valyala/gozstd"
	"github.com/zeebo/blake3"
)

func CompressWithGZip(src, dst string) {

	selectThreads := GetNumThreads()
	cLevel := GetGZipLevel()

	fsrc, fsrcInfo, fsrcHandler := NewBufReader(src)

	dstTemp := strings.Join([]string{dst, "ing"}, "")

	fdst, fh := NewBufWriter(dstTemp)

	w, err := gzip.NewWriterLevel(fdst, cLevel)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	var BufferMBByte int = BufferMB << 20

	w.SetConcurrency(BufferMBByte, selectThreads)

	log.Printf("threads: %v, buffer-size: %v MB", selectThreads, BufferMB)

	if isDebug {
		bar := pbar.NewBar64(fsrcInfo.Size())
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

func DecompressWithGZip(src string, dst string) error {
	fsrc, fsrcInfo, fhsrc := NewBufReader(src)

	dstTemp := strings.Join([]string{dst, "unzipping"}, ".")
	fdst, fhdst := NewBufWriter(dstTemp)

	reader, err := gzip.NewReader(fsrc)
	if err != nil {
		log.Fatal(err)
	}

	if isDebug {
		bar := pbar.NewBar64(fsrcInfo.Size())
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

	var buf []byte = make([]byte, BufferMB)
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

	var buf []byte = make([]byte, BufferMB)
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

	var buf []byte = make([]byte, BufferMB)
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

func setFilesMap(src string) (int64, error) {
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
	var fileSize int64 = 0
	var walkFunc = func(path string, info os.FileInfo, err error) error {
		path, _ = filepath.Abs(path)
		path = filepath.ToSlash(path)

		if !info.IsDir() {
			filesMap[path] = strings.Trim(strings.Replace(path, src[:strings.LastIndex(src, "/")], "", 1), "/")
			fileSize += info.Size()
		}

		return nil
	}
	err = filepath.Walk(src, walkFunc)
	return fileSize, err
}

func TarDir(src string, dst string) error {
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

	if isDebug {
		bar := pbar.NewBar64(0)
		err = format.Archive(context.Background(), io.MultiWriter(bufdst, bar), files)
		bar.Finish()
	} else {
		err = format.Archive(context.Background(), bufdst, files)
	}

	if err != nil {
		log.Fatal(err)
	}

	bufdst.Flush()
	fhdst.Close()

	err = os.Rename(dstTemp, dst)
	if err != nil {
		log.Fatal(err)
	}

	fullpathDst, _ := filepath.Abs(dst)
	Colorintln("green", "file: "+fullpathDst)

	return nil
}

func UntarDir(src string, dst string) error {
	fsrc, _, fhsrc := NewBufReader(src)
	wg := sync.WaitGroup{}

	format := archiver.CompressedArchive{
		Compression: nil,
		Archival:    archiver.Tar{},
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
		dstName = filepath.ToSlash(dstName)
		dstDir := filepath.ToSlash(filepath.Dir(dstName))
		if f.IsDir() {
			dstDir = filepath.ToSlash(dstName)
		}
		MakeDirs(dstDir)

		if f.IsDir() {
			return nil
		}

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

func PathNormalize(s string) string {
	var err error
	s, err = filepath.Abs(s)
	if err != nil {
		log.Fatal(err)
	}
	s = filepath.ToSlash(s)
	s = strings.TrimRight(s, "/")
	return s
}

func NewBufWriter(f string) (*bufio.Writer, *os.File) {
	fh, err := os.Create(f)
	if err != nil {
		log.Fatal(err)
	}

	return bufio.NewWriterSize(bufio.NewWriter(fh), BufferMB), fh
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

	return bufio.NewReaderSize(bufio.NewReader(fh), BufferMB), finfo, fh
}

func GetNumThreads() int {
	if Threads <= numCPU && Threads > 0 {
		return Threads
	}

	var autoThreads int = 1

	if numCPU > 1 && numCPU <= 4 {
		autoThreads = 2
	}

	if numCPU > 4 && numCPU <= 8 {
		autoThreads = 4
	}

	if numCPU > 8 {
		autoThreads = numCPU - 4
	}

	return autoThreads
}

func GetGZipLevel() int {
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

	return cLevel
}

func GetZstdLevel() zstd.EncoderLevel {
	cLevel := zstd.SpeedDefault
	switch Level {
	case 0:
		cLevel = zstd.SpeedFastest
	case 1:
		cLevel = zstd.SpeedDefault
	case 6:
		cLevel = zstd.SpeedBetterCompression
	case 9:
		cLevel = zstd.SpeedBestCompression
	default:
		cLevel = zstd.SpeedDefault
	}

	return cLevel
}

func CompressWithZstd(src, dst string) error {
	dstTemp := strings.Join([]string{dst, "ing"}, "")
	fdst, err := os.Create(dstTemp)
	if err != nil {
		return err
	}

	fsrc, fsrcInfo, fhsrc := NewBufReader(src)

	cLevel := GetZstdLevel()

	enc, err := zstd.NewWriter(fdst, zstd.WithEncoderLevel(cLevel))
	if Threads > 0 {
		numThreads := GetNumThreads()
		log.Println("threads:", numThreads)
		enc, err = zstd.NewWriter(fdst, zstd.WithEncoderLevel(cLevel), zstd.WithEncoderConcurrency(numThreads))
	}
	if err != nil {
		log.Fatal(err)
	}
	if isDebug {
		bar := pbar.NewBar64(fsrcInfo.Size())
		_, err = io.Copy(io.MultiWriter(enc, bar), fsrc)
		bar.Finish()
	} else {
		_, err = io.Copy(enc, fsrc)
	}

	if err != nil {
		enc.Close()
		log.Fatal(err)
	}
	enc.Close()

	fhsrc.Close()
	fdst.Close()

	err = os.Rename(dstTemp, dst)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func DecompressWithZstd(src, dst string) error {
	dstTemp := strings.Join([]string{dst, "ing"}, "")

	fdst, err := os.Create(dstTemp)
	if err != nil {
		log.Fatal(err)
	}

	fsrc, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}

	//dec, err = zstd.NewReader(fsrc, zstd.WithDecoderConcurrency(numThreads))
	dec, err := zstd.NewReader(fsrc)
	if err != nil {
		log.Fatal(err)
	}
	defer dec.Close()

	if isDebug {
		bar := pbar.NewBar64(0)
		_, err = io.Copy(io.MultiWriter(fdst, bar), dec)
		bar.Finish()
	} else {
		_, err = io.Copy(fdst, dec)
	}

	if err != nil && err != io.EOF {

		log.Println("error: io.copy")
		log.Fatal(err)
	}

	fsrc.Close()

	fdst.Close()

	err = os.Rename(dstTemp, dst)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
