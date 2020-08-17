package samutils

import "testing"

//32位 小写
//@author sam@2020-08-17 13:43:57
func TestMd5(t *testing.T) {
	str:="sam"
	t.Log(Md5(str))
	t.Log(Md5Check(str,"332532dcfaa1cbf61e2a266bd723612c"))
}

//测试密码生成
//@author sam@2020-08-17 13:44:47
func  TestEncryptPassword(t *testing.T){
	t.Log(EncryptPassword("samiaagoodman","cleEM`\"p"))
}
