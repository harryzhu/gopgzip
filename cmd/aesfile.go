package cmd

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"io"
	"strings"

	"encoding/hex"

	"log"
	"os"

	progressbar "github.com/schollz/progressbar/v3"
)

const AESCHUNKSIZE int64 = 16 << 20
const AESBLOCKSIZE int = 16

// PwdKey length can be 16

var (
	PwdKey []byte
	IVKey  []byte
)

// ------------

func setKeyPasswordIV() {
	var salt string = SHA256("Cu5t0m-s@lt")
	var p string
	if Password != "" && Force {
		p = Password
	} else {
		p = GetEnv("HARRYZHUENCRYPTKEY", passwordDefault)
	}

	if p == "" {
		log.Fatal("you did not set any password")
	}

	pk := SHA256(MD5(p) + ":" + salt)
	ivk := SHA256(MD5(pk) + ":" + salt)

	PwdKey = []byte(pk)[:AESBLOCKSIZE]
	IVKey = []byte(ivk)[:AESBLOCKSIZE]

	if PwdKey == nil || IVKey == nil {
		log.Fatal("password and iv key cannot be empty")
	}
}

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

	Colorintln("green", "file: "+dst)
	return nil
}

func GetEnv(s string, vDefault string) string {
	v := os.Getenv(s)
	if v == "" {
		return vDefault
	}
	return strings.Trim(v, " ")
}

func MD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func SHA256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
