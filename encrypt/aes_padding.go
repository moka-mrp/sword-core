package encrypt

import "bytes"

// 填充 缺几位补几位 且填充的位数和填充的数字相同
// PKCS7Padding 填充 块大小可指定
func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// PKCS5Padding 填充 块大小固定为8
func PKCS5Padding(ciphertext []byte) []byte {
	return PKCS7Padding(ciphertext, 8)
}

// UnPadding 移除填充
func UnPadding(data []byte) []byte {
	// 获取最后一位的数字值 unpadding
	length := len(data)
	unpadding := int(data[length-1])
	// 移除 unpadding个字符
	return data[:(length - unpadding)]
}
