package aliots

import "github.com/moka-mrp/sword-core/config"

var conf config.OtsMultiConfig
func init()  {
	conf=make(config.OtsMultiConfig,2)
	//conf["tulong"]=config.OtsConfig{
	//	EndPoint:        "",
	//	InstanceName:    "",
	//	AccessKeyId:     "",
	//	AccessKeySecret: "",
	//}
	//
	//
	//conf["xueyin"]=config.OtsConfig{
	//	EndPoint:        "",
	//	InstanceName:    "",
	//	AccessKeyId:     "",
	//	AccessKeySecret: "",
	//}

	conf["tulong"]=config.OtsConfig{
		EndPoint:        "https://tulong.cn-shanghai.ots.aliyuncs.com",
		InstanceName:    "tulong",
		AccessKeyId:     "LTAI4GG1DPRi2yty32FqBdj3",
		AccessKeySecret: "FUpx9HHvJ2SUiSuVk6OaI4SehsGM3A",
	}


	conf["xueyin"]=config.OtsConfig{
		EndPoint:        "https://xueyin.cn-shanghai.ots.aliyuncs.com",
		InstanceName:    "xueyin",
		AccessKeyId:     "LTAI4GG1DPRi2yty32FqBdj3",
		AccessKeySecret: "FUpx9HHvJ2SUiSuVk6OaI4SehsGM3A",
	}






}