// Package hash 哈希算法
// sha1 sha2(256)算法
package hash

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// ByteSha1 计算Sha1
func ByteSha1(data []byte) string {
	hashInstance := sha1.New()
	hashInstance.Write(data)
	return hex.EncodeToString(hashInstance.Sum([]byte("")))
}

// Sha1 计算Sha1
func Sha1(str string) string {
	hashInstance := sha1.New()
	hashInstance.Write([]byte(str))
	return hex.EncodeToString(hashInstance.Sum([]byte("")))
}

// FileSha1 计算文件sha1
func FileSha1(file *os.File) string {
	hashInstance := sha1.New()
	io.Copy(hashInstance, file)
	return hex.EncodeToString(hashInstance.Sum(nil))
}

// ByteSha256 计算Sha256
func ByteSha256(data []byte) string {
	hashInstance := sha256.New()
	hashInstance.Write(data)
	return hex.EncodeToString(hashInstance.Sum([]byte("")))
}

// Sha256 计算Sha256
func Sha256(str string) string {
	hashInstance := sha256.New()
	hashInstance.Write([]byte(str))
	return hex.EncodeToString(hashInstance.Sum([]byte("")))
}

// Sha1Check 校验Sha1
func Sha1Check(str string, sha1 string) bool {
	return sha1 == Sha1(str)
}

// Sha256Check 校验Sha256
func Sha256Check(str string, sha256 string) bool {
	return sha256 == Sha256(str)
}
