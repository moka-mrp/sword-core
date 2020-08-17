package samutils

import (
	"fmt"
	"math/rand"
	"time"
)


// 定义字符串
var (
	//特殊字符
	stringMisc = []byte(".$#@&*_-")
	//随机阿拉伯数字
	stringDigit    = []byte("1234567890")
	stringDigitLen = len(stringDigit)
	//小写的26字符随机
	stringLWord    = []byte("abcdefghijklmnopqrstuvwxyz")
	stringLWordLen = len(stringLWord)
    //大写的26字符随机
	stringUWord    = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	stringUWordLen = len(stringUWord)
	//大小写的52字符随机
	stringCWord    = []byte(fmt.Sprintf("%s%s", stringLWord, stringUWord))
	stringCWordLen = len(stringCWord)
	//大写、小写、数字
	stringL3    = []byte(fmt.Sprintf("%s%s%s", stringDigit, stringLWord,stringUWord))
	stringL3Len = len(stringL3)
    //大写、小写、数字、特殊字符的随机
	stringL4    = []byte(fmt.Sprintf("%s%s%s%s", stringDigit, stringLWord, stringMisc, stringUWord))
	stringL4Len = len(stringL4)
)



//随机阿拉伯数字
// Intn以int形式返回[0，n）中的非负伪随机数,如果n <= 0，则发生恐慌。
// @author sam@2020-08-17 11:19:32
func RandStringDigit(n int64) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = stringDigit[rand.Intn(stringDigitLen)]
	}
	return string(b)
}

//小写的26字符随机
//@author sam@2020-08-17 11:20:55
func RandStringWordL(n int64) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = stringLWord[rand.Intn(stringLWordLen)]
	}
	return string(b)
}

//大写的26字符随机
//@author sam@2020-08-17 11:21:05
func RandStringWordU(n int64) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = stringUWord[rand.Intn(stringUWordLen)]
	}
	return string(b)
}

//大小写的52字符随机
//@author sam@2020-08-17 11:21:20
func RandStringWordC(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = stringCWord[rand.Intn(stringCWordLen)]
	}
	return string(b)
}




//大写、小写、数字
//@author sam@2020-08-17 11:30:37
func RandStringL3(n int64) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = stringL3[rand.Intn(stringL3Len)]
	}
	return string(b)
}


//大写、小写、数字、特殊字符的随机
//@author sam@2020-08-17 11:30:37
func RandStringL4(n int64) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = stringL4[rand.Intn(stringL4Len)]
	}
	return string(b)
}
//--------------------------------------------------------------------------------------------------------------------

//获取盐串
//@author sam@2020-08-17 11:47:34
func GenSalt() string {
	return RandStringL3(8)
}




//产生min-max中的一个随机数，[min,max)
//@author sam@2020-08-17 11:44:00
func MtRand(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if min > max {
		return max
	}
	//[0,max-min) +min   [min,max)
	return r.Intn(max-min) + min //// Intn以int形式返回[0，n）中的非负伪随机数
}

