package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	tempStr := h.Sum(nil)
	return hex.EncodeToString(tempStr)
}
func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode(data))
}

// 加密
func MakePassword(plainpwd, salt string) string {
	return Md5Encode(plainpwd + salt)
}
func ValidPassword(plainpwd, salt string, password string) bool {
	md := MakePassword(plainpwd, salt)
	fmt.Println(md + "				" + password)
	return md == password
}
