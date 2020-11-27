package aliots

import "github.com/moka-mrp/sword-core/config"

var conf config.OtsMultiConfig
func init()  {
	conf=make(config.OtsMultiConfig,2)
	conf["tulong"]=config.OtsConfig{
		EndPoint:        "",
		InstanceName:    "",
		AccessKeyId:     "",
		AccessKeySecret: "",
	}


	conf["xueyin"]=config.OtsConfig{
		EndPoint:        "",
		InstanceName:    "",
		AccessKeyId:     "",
		AccessKeySecret: "",
	}








}