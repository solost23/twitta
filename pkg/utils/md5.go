package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

func NewMd5(str string, salt ...interface{}) string {
	if len(salt) > 0 {
		slice := make([]string, len(salt)+1)
		str = fmt.Sprintf(str+strings.Join(slice, "%v"), salt...)
	}
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func GetMd5Hash(str string) string {
	sum := md5.Sum([]byte(str))
	return hex.EncodeToString(sum[:])
}
