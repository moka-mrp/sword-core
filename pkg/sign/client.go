package sign

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

type Client struct {
	jWay               string  //拼接方式
	snType             string   // 签名类型
	signatureExpired int64      // 签名有效期
	sk    string      //签名秘钥
	p     Params      //参与签名的参数
}

// 创建验签客户端
// @author sam@2020-06-18 10:07:08
func NewClient(params Params,sk string) *Client {
	return &Client{
		jWay:JwUrl,
		snType:         MD5,
		signatureExpired: SignatureExpired,
		sk:               sk,
		p:params,
	}
}

// 验证签名
func (c *Client) ValidSign() (b bool,err error) {
	//必传的验签参数检测
	if len(c.p.GetString(Sn))<=0  || len(c.p.GetString(Timestamp))<=0 || len(c.p.GetString(Nonce))<=0 {
		err=errors.New("params missing")
		return
	}
	//判断签名是否失效
	timestamp:=c.p.GetInt64("timestamp")
	if time.Now().Unix()  - timestamp > c.signatureExpired {
	  return false,errors.New("signature expired")
	}
	//动态调整参数
	if len(c.p.GetString("sn_type")) >0{
		c.snType=c.p.GetString("sn_type")
	}
	//验证签名
	//fmt.Println("服务端计算的签名为",c.Sign())
	b=(c.p.GetString(Sn) == c.Sign())
	if !b{
		err=errors.New("sign failed")
	}
	return
}


//向请求参数中添加   sign_type、sn 必填参数
//@author sam@2020-06-18 10:48:04
func (c *Client) fillRequestData(params Params) Params {
	params["nonce"] = nonceStr()
	params["timestamp"]=string(time.Now().Unix())
	params["sn_type"] = c.snType
	params["sn"] = c.Sign()
	return params
}



//生成签名
//@todo  如果参数的值为空不参与签名
//@todo  sn本身不参与签名,所以不要使用该关键字额
//@author sam@2020-06-18 10:55:32
func (c *Client) Sign() string {
	//1.拼接参数
	buf:=c.JoinWay()
	//fmt.Println("拼接为:",string(buf.Bytes()))
	//2.签名运算
	var (
		dataMd5    [16]byte
		dataSha256 []byte
		dataSha1 []byte
		str        string
	)
	switch c.snType {
	case MD5:
		dataMd5 = md5.Sum(buf.Bytes())
		str = hex.EncodeToString(dataMd5[:]) //需转换成切片
	case HMACSHA256:
		h := hmac.New(sha256.New, []byte(c.sk))
		h.Write(buf.Bytes())
		dataSha256 = h.Sum(nil)
		str = hex.EncodeToString(dataSha256[:])
	case SHA1:
		h := sha1.New()
		h.Write(buf.Bytes())
		dataSha1=h.Sum(nil)
		str=fmt.Sprintf("%x",dataSha1)
	}
	//3.再将得到的字符串所有字符转换为大写，得到sign值signValue。
	return strings.ToUpper(str)
}

//应付客户端传递的不同的拼接方式
//todo 这里后续可以拓展成不同种类的拼接方式额
//@author sam@2020-11-04 10:48:50
func (c *Client)JoinWay()(buf bytes.Buffer){
	//设所有发送或者接收到的参数key放到切片keys中
	var keys = make([]string, 0, len(c.p))
	//遍历签名参数,过滤sign以及为空的参数
	for k := range c.p {
		if k != "sn" { // 排除sn
			keys = append(keys, k)
		}
	}

	switch c.jWay{
	case JwUrl:
		sort.Strings(keys) // 参数名ASCII码从小到大排序（字典序)
		// 使用URL键值对的格式（即key1=value1&key2=value2…）拼接成字符串A
		//var buf bytes.Buffer //创建字符缓冲,用来快速拼接字符串
		for _, k := range keys {
			if len(c.p.GetString(k)) > 0 { //如果参数的值为空不参与签名
				buf.WriteString(k)
				buf.WriteString(`=`)
				buf.WriteString(c.p.GetString(k))
				buf.WriteString(`&`)
			}
		}
		//在字符串A最后拼接上sk得到B字符串
		buf.WriteString(`sk=`)
		buf.WriteString(c.sk)
	}

	return
}

//动态设置签名过期时间
//@author sam@2020-11-04 10:52:22
func (c *Client) SetSignatureExpired(duration int64) {
	c.signatureExpired=duration
}

//动态设置签名类型
//@author sam@2020-11-04 10:52:14
func (c *Client) SetSignType(signType string) {
	c.snType = signType
}

//动态设置字符串拼接方式
//@author sam@2020-11-04 10:52:05
func (c *Client) SetJway(jWay string) {
	c.jWay = jWay
}

//动态设置签名秘钥
//@author sam@2020-11-04 10:51:55
func (c *Client) SetSk(sk string) {
	c.sk=sk
}
//动态设置签名参数
//@author sam@2020-11-04 10:51:48
func (c *Client) SetParams(p Params) {
	c.p=p
}

