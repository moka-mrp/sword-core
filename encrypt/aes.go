// Package encrypt 加密 aes对称加密
package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
)

// AesEncrypt Aes加密
func AesEncrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key) // 得到一些数字
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()                    // 得到16
	originBytes := PKCS7Padding(plaintext, blockSize) // 得到补位结果

	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) // cbc算法初始化，抛弃key的一半，当做向量

	chipher := make([]byte, len(originBytes))   // 创建一个空容器
	blockMode.CryptBlocks(chipher, originBytes) // aes加密结果 参数1 空容器，参数2，源数据
	return chipher, nil
}

// AesDecrypt Aes解密
func AesDecrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key) // 得到一些数字
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()                              // 得到16
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) // cbc算法初始化，抛弃key的一半，当做向量

	originBytes := make([]byte, len(ciphertext))   // 创建空容器
	blockMode.CryptBlocks(originBytes, ciphertext) // aes解密结果

	originBytes = UnPadding(originBytes) // 移除补位

	return originBytes, nil
}

// AesEncryptString 加密
func AesEncryptString(plaintext string, key []byte) (string, error) {
	ciphertext, err := AesEncrypt([]byte(plaintext), key)
	if err != nil {
		return "", err
	}
	return ByteBase64Encode(ciphertext), nil
}

// AesDecryptString 解密
func AesDecryptString(ciphertext string, key []byte) (string, error) {
	ciphertext = Base64Decode(ciphertext)
	cipherByte, err := AesDecrypt([]byte(ciphertext), key)
	if err != nil {
		return "", err
	}
	return string(cipherByte), nil
}
