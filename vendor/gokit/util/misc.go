package util

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"os"
	"path"
)

func GetRoot() (root string) {
	cmd := os.Args[0]
	if path.IsAbs(cmd) {
		root = path.Dir(cmd)
	} else {
		wd, _ := os.Getwd()
		root = path.Dir(path.Join(wd, cmd))
	}
	return
}

func GetFile(file string) string {
	if path.IsAbs(file) {
		return file
	}

	return path.Join(GetRoot(), file)
}

func RandomString(length uint64) (string, error) {
	bl := length / 2
	if length%2 != 0 {
		bl++
	}

	buf := make([]byte, bl)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(buf), nil
}
