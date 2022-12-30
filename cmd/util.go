package cmd

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"strings"
	"sync"

	"github.com/klauspost/compress/zstd"
	gzip "github.com/klauspost/pgzip"

	"github.com/zeebo/blake3"
	"github.com/zeebo/xxh3"
)

func CompressWithGZip(src, dst string) {
	selectThreads := GetNumThreads()
	cLevel := GetGZipLevel()

	fsrc, fsrcInfo, fhsrc := NewBufReader(src)

	dstTemp := strings.Join([]string{dst, "ing"}, "")

	fdst, fh := NewBufWriter(dstTemp)

	w, err := gzip.NewWriterLevel(fdst, cLevel)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	var BufferMBByte int = BufferMB << 20
	w.Name = fsrcInfo.Name()
	w.SetConcurrency(BufferMBByte, selectThreads)

	log.Printf("threads: %v, buffer-size: %v MB", selectThreads, BufferMB)

	bar64.WithMax64(fsrcInfo.Size())
	_, err = io.Copy(io.MultiWriter(w, bar64), fsrc)
	bar64.Finish()

	if err != nil {
		log.Fatal(err)
	}

	w.Close()
	fdst.Flush()
	fh.Close()
	fhsrc.Close()

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

	bar64.WithMax64(fsrcInfo.Size())
	_, err = reader.WriteTo(io.MultiWriter(fdst, bar64))
	fdst.Flush()
	bar64.Finish()

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
	fsrc, _, fhsrc := NewBufReader(src)
	hash := md5.New()

	var buf []byte = make([]byte, 8192)
	for {
		n, err := fsrc.Read(buf)
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

	var buf []byte = make([]byte, 8192)
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

	var buf []byte = make([]byte, 8192)
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

func LoadExcludes() error {
	if Excludes == "" {
		return nil
	}

	Excludes, _ = filepath.Abs(Excludes)
	f, err := os.Open(Excludes)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("excludes file was loaded:", Excludes)
	}

	c, _ := ioutil.ReadAll(f)
	content := string(c)
	content = strings.ReplaceAll(content, "\\r\\n", "\\n")
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		ExcludeFiles = append(ExcludeFiles, line)
	}

	return nil
}

func IsExcluded(s string) bool {
	if len(ExcludeFiles) == 0 {
		return false
	}

	s = strings.TrimSpace(s)

	if s == "" {
		return false
	}

	s = filepath.ToSlash(s)
	for _, line := range ExcludeFiles {
		line = filepath.ToSlash(line)

		if line == s {
			return true
		}

		if strings.HasPrefix(s, line) {
			return true
		}

		if strings.HasPrefix(filepath.Base(s), line) {
			return true
		}

	}

	return false
}

func setFilesMap(src string) (int64, error) {
	LoadExcludes()

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

		nameInTar := strings.Trim(strings.Replace(path, src[:strings.LastIndex(src, "/")], "", 1), "/")

		bExc := IsExcluded(nameInTar)
		if bExc == true {
			log.Println("EXCLUDE:", nameInTar)
			return nil
		}

		if !info.IsDir() {
			filesMap[path] = nameInTar
			fileSize += info.Size()
		}

		return nil
	}
	err = filepath.Walk(src, walkFunc)
	return fileSize, err
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
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	if strings.Contains(s, ":") && strings.Index(s, ":") < 10 {
		return s
	}

	if strings.HasPrefix(s, "\\\\") {
		return s
	}

	s = filepath.ToSlash(s)
	var err error

	s, err = filepath.Abs(s)
	if err != nil {
		log.Fatal(err)
	}

	s = strings.TrimRight(s, "/")
	s = FormatString(s)
	if isDebug {
		log.Println("PathNormalize:", s)
	}
	return s
}

func Filepathify(fp string) string {
	var replacement string = "_"

	reControlCharsRegex := regexp.MustCompile("[\u0000-\u001f\u0080-\u009f]")

	reRelativePathRegex := regexp.MustCompile(`^\.+`)

	filenameReservedRegex := regexp.MustCompile(`[<>:"\\|?*\x00-\x1F]`)
	filenameReservedWindowsNamesRegex := regexp.MustCompile(`(?i)^(con|prn|aux|nul|com[0-9]|lpt[0-9])$`)

	// reserved word
	fp = filenameReservedRegex.ReplaceAllString(fp, replacement)

	// continue
	fp = reControlCharsRegex.ReplaceAllString(fp, replacement)
	fp = reRelativePathRegex.ReplaceAllString(fp, replacement)
	fp = filenameReservedWindowsNamesRegex.ReplaceAllString(fp, replacement)
	return fp
}

func NewBufWriter(f string) (*bufio.Writer, *os.File) {
	fh, err := os.Create(f)
	if err != nil {
		fh.Close()
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
		fh.Close()
		log.Fatal(err)
	}

	return bufio.NewReaderSize(bufio.NewReader(fh), BufferMB), finfo, fh
}

func GetNumThreads() int {
	if Threads <= NumCPU && Threads > 0 {
		return Threads
	}

	var autoThreads int = 1

	if NumCPU > 1 && NumCPU <= 4 {
		autoThreads = 2
	}

	if NumCPU > 4 && NumCPU <= 8 {
		autoThreads = 4
	}

	if NumCPU > 8 {
		autoThreads = NumCPU - 4
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

	bar64.WithMax64(fsrcInfo.Size())
	_, err = io.Copy(io.MultiWriter(enc, bar64), fsrc)
	bar64.Finish()

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

	bar64.WithMax64(0)
	_, err = io.Copy(io.MultiWriter(fdst, bar64), dec)
	bar64.Finish()

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

func RatioInputOutput(src string, dst string) error {
	_, fsrcInfo, _ := NewBufReader(src)
	_, fdstInfo, _ := NewBufReader(dst)
	r := float64(fdstInfo.Size()) / float64(fsrcInfo.Size())
	log.Println("compress ratio:", r)
	return nil
}

func Xxh3SumFile(src string) string {
	hash := xxh3.New()
	fsrc, _, fhsrc := NewBufReader(src)

	var buf []byte = make([]byte, 8192)
	for {
		n, err := fsrc.Read(buf)
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

func Xxh3String(src string) string {
	hash := xxh3.New()
	hash.Write([]byte(src))
	return hex.EncodeToString(hash.Sum(nil))
}

func FileInfo(dst string) (finfo fs.FileInfo, err error) {
	finfo, err = os.Stat(dst)
	if err != nil {
		return nil, err
	}
	return finfo, nil
}

func CopyFile(src string, dst string) error {
	_, err := FileInfo(dst)
	if err == nil && IsOverwrite == false {
		if isDebug {
			log.Println("file exists and will not overwrite:", dst)
		}
		return nil
	}

	_, err = FileInfo(src)
	if err != nil {
		if isDebug {
			log.Println("file source does not exist:", src)
		}
		return err
	}

	MakeDirs(filepath.Dir(dst))

	fsrc, fsrcInfo, fhsrc := NewBufReader(src)
	fdst, fhdst := NewBufWriter(dst)

	bar64.WithMax64(fsrcInfo.Size())
	_, err = io.Copy(io.MultiWriter(fdst, bar64), fsrc)
	fdst.Flush()
	bar64.Finish()

	fhsrc.Close()
	fhdst.Close()

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func CopyDir(src string, dst string) error {
	var copyList map[string]string = make(map[string]string, 8192)
	var copySum int

	var dstPath string
	var fullPath string
	var err error

	var walkFunc = func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		fullPath, err = filepath.Abs(path)
		if err != nil {
			log.Println(err)
			return err
		}
		fullPath = filepath.ToSlash(path)

		dstPath = strings.ReplaceAll(fullPath, filepath.Dir(src), dst)
		dstPath = filepath.ToSlash(dstPath)
		dstPath = strings.TrimRight(dstPath, "/")

		if IsOverwrite == false {
			_, err := FileInfo(dstPath)
			if err == nil {
				return nil
			}
		}

		if fullPath == "" || dstPath == "" {
			return nil
		}

		copyList[fullPath] = dstPath
		copySum += 1

		return nil
	}

	err = filepath.Walk(src, walkFunc)

	if err != nil {
		log.Fatal("filepath.Walk:", err)
	}

	bar64.WithDisabled64(true)

	bar.WithMax(copySum)
	bar.WithCounterSkip(100)
	bar.WithCounterCycle(20)

	for fsrc, fdst := range copyList {
		bar.Add(1)
		CopyFile(fsrc, fdst)
	}

	bar.Finish()

	return nil
}

func FormatString(dst string) string {
	if dst == "" {
		return ""
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "-"
	} else {
		hostname = strings.ToLower(hostname)
	}

	yyyy := tStart.Format("2006")
	mm := tStart.Format("01")
	dd := tStart.Format("02")
	HH := tStart.Format("15")
	MM := tStart.Format("04")
	SS := tStart.Format("05")

	DayOfWeek := strings.ToLower(tStart.Weekday().String())

	dst = strings.ReplaceAll(dst, "{hostname}", hostname)
	dst = strings.ReplaceAll(dst, "{yyyy}", yyyy)
	dst = strings.ReplaceAll(dst, "{mm}", mm)
	dst = strings.ReplaceAll(dst, "{dd}", dd)
	dst = strings.ReplaceAll(dst, "{HH}", HH)
	dst = strings.ReplaceAll(dst, "{MM}", MM)
	dst = strings.ReplaceAll(dst, "{SS}", SS)
	dst = strings.ReplaceAll(dst, "{day-of-week}", DayOfWeek)

	return dst
}

func DownloadFile(src string, dst string) error {
	resp, err := http.Get(src)
	if err != nil {
		log.Println("Error(http.Get):", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 {
		log.Println(resp.StatusCode, "(ERROR: cannot get url):", src)
		return errors.New("cannot get the url: " + src)
	}

	//
	dst = PathNormalize(dst)

	_, err = os.Stat(dst)
	if err == nil {
		if IsOverwrite == false {
			log.Println("SKIP(as exists):", dst)
			return nil
		} else {
			err = os.Remove(dst)
			if err != nil {
				log.Println("Error(os.Remove):", err)
			}
		}
	}
	if isDebug {
		log.Println("start downloading:", src, "==>", dst)
	}

	dstTemp := strings.Join([]string{dst, "downloading"}, ".")
	MakeDirs(filepath.Dir(dstTemp))
	fdst, fhdst := NewBufWriter(dstTemp)
	defer fhdst.Close()

	bar64.WithMax64(resp.ContentLength)
	_, err = io.Copy(io.MultiWriter(fdst, bar64), resp.Body)
	fdst.Flush()
	bar64.Finish()

	if err != nil {
		log.Println("Error(io.Copy):", err)
		return err
	}

	fhdst.Close()

	err = os.Rename(dstTemp, dst)
	if err != nil {
		log.Println("Error(os.Remove):", err)
		return err
	}

	return nil
}

func DownloadByList(src string, dstDir string) error {
	bsrc, err := ioutil.ReadFile(src)
	if err != nil {
		log.Println(err)
		return err
	}

	strSrc := string(bsrc)
	strSrc = strings.ReplaceAll(strSrc, "\r\n", "\n")
	srcLines := strings.Split(strSrc, "\n")

	var downList []string

	for _, line := range srcLines {
		line = strings.Trim(line, "\\n")
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		downList = append(downList, line)
	}

	wg := sync.WaitGroup{}
	var downCount int

	for _, line := range downList {
		dstPath := ""
		srcUrl := ""
		if strings.Index(line, "|") == -1 {
			srcUrl = strings.TrimSpace(line)
			linePath := line[strings.Index(line, "://")+3:]
			linePath = strings.Trim(linePath, "/")
			linePath = Filepathify(linePath)
			dstPath = filepath.Join(dstDir, linePath)
			if IsKeepUrlPath == false {
				dstPath = filepath.Join(dstDir, filepath.Base(linePath))
			}
		} else {
			srcDst := strings.Split(line, "|")
			if len(srcDst) != 2 {
				log.Println("invalid line:", line)
				log.Println("if you are using `|`, please be sure the format is `src_remote_url|local_file_save_path`")
				continue
			}
			srcUrl = strings.TrimSpace(srcDst[0])
			dstPath = strings.TrimSpace(srcDst[1])

		}

		dstPath = filepath.ToSlash(dstPath)

		wg.Add(1)
		downCount += 1
		go func() {
			DownloadFile(srcUrl, dstPath)
			downCount -= 1
			wg.Done()
		}()

		if downCount >= 10 {
			wg.Wait()
		}

	}

	wg.Wait()

	return nil
}
