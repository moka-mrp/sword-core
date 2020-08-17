package samutils

import "testing"


//测试随机数字
//@author sam@2020-08-17 11:16:57
func TestRandStringDigit(t *testing.T) {
	t.Log(RandStringDigit(5))
}

//测试随机小写字母
//@author sam@2020-08-17 11:20:01
func TestRandStringWordL(t *testing.T) {
	t.Log(RandStringWordL(5))
}

//测试随机大写字母
//@author sam@2020-08-17 11:22:19
func TestRandStringWordU(t *testing.T) {
	t.Log(RandStringWordU(5))
}

//测试随机大小写字母
//@author sam@2020-08-17 11:22:47
func TestRandStringWordC(t *testing.T) {
	t.Log(RandStringWordC(5))
}

//测试 大写、小写、数字
//@author sam@2020-08-17 11:37:44
func TestRandStringL3(t *testing.T) {
	t.Log(RandStringL3(5))
}


//测试 大写、小写、数字、特殊字符
//@author sam@2020-08-17 11:38:25
func TestRandStringL4(t *testing.T) {
	t.Log(RandStringL4(5))
}
//--------------------------------------

//测试随机区间数的返回[0,3)
//@author sam@2020-08-17 11:45:53
func TestMtRand(t *testing.T) {
	t.Log(MtRand(0,3))
}

//测试盐串的产生
//@author sam@2020-08-17 11:48:34
func TestGenSalt(t *testing.T) {
	t.Log(GenSalt())
}

