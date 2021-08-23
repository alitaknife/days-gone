package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/gogf/gf/net/ghttp"
	"io"
)

func Sha1Encrypt(file *ghttp.UploadFile) (s string, err error) {
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	h := sha1.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), err
}
