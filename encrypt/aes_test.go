package encrypt

import (
	"log"
	"testing"
)

func TestAes(t *testing.T) {
	// 定义加密用的key
	keyStr := "abcdefghijklmnopqrstuvwxyz123456" // 32个字符串
	key := []byte(keyStr)                        // 2个字符串对应一个 byte字节 也就是16字节  每个字节占8位 16*8=128 对应Aes-128

	// 加密
	plaintext := "go is perfect langage"
	cipherByte, _ := AesEncrypt([]byte(plaintext), key)
	log.Println(cipherByte)
	cipherStr, _ := AesEncryptString(plaintext, key)
	log.Println(cipherStr)

	// 解密
	log.Println("解密")
	oldByte, _ := AesDecrypt([]byte(cipherByte), key)
	log.Println(string(oldByte))

	oldStr, _ := AesDecryptString(cipherStr, key)
	log.Println(oldStr)

}
