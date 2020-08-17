package samutils

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5加密
func Md5(s string) string {
	md := md5.New()
	md.Write([]byte(s))
	return hex.EncodeToString(md.Sum(nil))
}

//检查MD5加密
func Md5Check(s, m string) bool {
	return m == Md5(s)
}



//明文密码结合盐串进行加密而已
//@todo 注意这里并不是单纯的md5加密
//@author sam@2020-08-17 11:51:59
func EncryptPassword(pwd, salt string) string {
	//return Md5(pwd+salt)
	m := md5.Sum([]byte(pwd + salt))
	m = md5.Sum(m[:])
	return hex.EncodeToString(m[:])
}


