package cmd

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"io"
	"time"

	"encoding/hex"

	"log"
	"os"

	progressbar "github.com/schollz/progressbar/v3"
)

const AESCHUNKSIZE int64 = 16 << 20
const AESBLOCKSIZE int = 16

// PwdKey length can be 16

var PwdKey = []byte(MD5(GetEnv("HAZHUENCRYPTKEY", "This(*Key*)@2021This(*Key*)@2021")))[:AESBLOCKSIZE]
var IVKey = []byte(MD5(GetEnv("HAZHUENCRYPTKEY", "That(*Key*)@2021That(*Key*)@2021")))[:AESBLOCKSIZE]

// ------------

func AESEncodeFile(src string, dst string) error {
	fsrc, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer fsrc.Close()

	fdst, err := os.Create(dst)
	if err != nil {
		log.Fatal(err)
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		log.Fatal(err)
	}

	tStart := time.Now()
	bar := progressbar.DefaultBytes(srcInfo.Size())

	iv := []byte(IVKey)

	block, err := aes.NewCipher(PwdKey)
	if err != nil {
		log.Fatal(err)
	}

	stream := cipher.NewCTR(block, iv)

	srcReader := bufio.NewReader(fsrc)
	buf := make([]byte, AESCHUNKSIZE)

	for {
		n, err := srcReader.Read(buf)
		if n == 0 {
			if err == io.EOF {
				//log.Println("EOF")
				break
			}

			if err != nil {
				log.Println(err)
				break
			}

		}
		encByte := make([]byte, n)
		stream.XORKeyStream(encByte, buf[:n])

		_, err = fdst.Write(encByte)
		if err != nil {
			log.Fatal(err)
		}
		bar.Add64(int64(n))
	}

	bar.Finish()
	tStop := time.Now()
	duration := tStop.Sub(tStart)

	log.Printf("OK. duration: %v sec\n", duration)

	Colorintln("green", "file: "+dst)
	return nil
}

func AESDecodeFile(src string, dst string) error {
	fsrc, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer fsrc.Close()

	fdst, err := os.Create(dst)
	if err != nil {
		log.Fatal(err)
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		log.Fatal(err)
	}

	tStart := time.Now()

	bar := progressbar.DefaultBytes(srcInfo.Size())

	iv := []byte(IVKey)

	block, err := aes.NewCipher(PwdKey)
	if err != nil {
		log.Fatal(err)
	}

	stream := cipher.NewCTR(block, iv)

	srcReader := bufio.NewReader(fsrc)
	buf := make([]byte, AESCHUNKSIZE)

	for {
		n, err := srcReader.Read(buf)

		if n == 0 {

			if err == io.EOF {
				//log.Println("EOF")
				break
			}

			if err != nil {
				log.Println(err)
				break
			}

		}
		decByte := make([]byte, n)
		stream.XORKeyStream(decByte, buf[:n])

		_, err = fdst.Write(decByte)
		if err != nil {
			log.Fatal(err)
		}
		bar.Add64(int64(n))
	}
	bar.Finish()
	tStop := time.Now()
	duration := tStop.Sub(tStart)

	log.Printf("OK. duration: %v sec\n", duration)
	Colorintln("green", "file: "+dst)
	return nil
}

func GetEnv(s string, vDefault string) string {
	v := os.Getenv(s)
	if v == "" {
		return vDefault
	}
	return v
}

func MD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
