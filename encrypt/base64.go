// Package encrypt base64 编码解码
package encrypt

// base64是简单的移位替换密码
// base58是base64的子集,需要手动实现
import (
	"encoding/base64"
)

// Base64Encode 编码
func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// ByteBase64Encode 编码
func ByteBase64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Base64Decode 解码
func Base64Decode(str string) string {
	res, _ := base64.StdEncoding.DecodeString(str)
	return string(res)
}
