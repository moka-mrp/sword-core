package sign

import (
	"net/http"
	"strconv"
	"time"
)



//gin的get参数转map
//todo  url.Values{"ak":[]string{"11111"}, "name":[]string{"sam2222", "sam"}}
//todo 比如get传递了name=sam，只有post按照x-www-form-urlencoded方式才会获取第二个name值，比如sam2222
//@author sam@2020-11-03 15:53:00
func UrlencodedParamsToMap(r *http.Request) Params {
    //解析请求传递的参数(get的url以及post的参数)
	r.ParseForm()
	//这里我们只追对urlencoded编码的参数进行分析
	params := make(Params)
	for k,v:=range r.Form{
		//fmt.Println(k,v[0])
		params.SetString(k,v[0])
	}
	return params
}


// 用时间戳生成随机字符串
// todo 故意取了UTC时区额
// @author sam@2020-06-18 10:48:53
func nonceStr() string {
	return strconv.FormatInt(time.Now().UTC().UnixNano(), 10) //1592448742386052000
}


