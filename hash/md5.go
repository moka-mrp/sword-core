package hash

import (
	"crypto/md5"
	"encoding/hex"
)

// Md5 快速计算Md5
func Md5(str string) string {
	hashInstance := md5.New()
	hashInstance.Write([]byte(str))
	return hex.EncodeToString(hashInstance.Sum(nil))
}

// Md5Check 校验Md5
func Md5Check(str string, md5 string) bool {
	return md5 == Md5(str)
}

// ByteMd5 传递[]byte计算md5
func ByteMd5(data []byte) string {
	hashInstance := md5.New()
	hashInstance.Write(data)
	return hex.EncodeToString(hashInstance.Sum(nil))
}
