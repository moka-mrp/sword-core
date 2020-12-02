package samutils

import "time"

//获取秒数
//@author sam@2020-11-30 11:56:19
func GetSeconds(dur time.Duration) int64 {
	sec := dur / time.Second
	return int64(sec)
}


//获取时间统计前缀
//@author sam@2020-12-02 10:15:45
func GetPrefixTime(st int64, timeType string) string {
	var layout string
	switch timeType {
	case "D":
		layout = "20060102"
		return "D-"+time.Unix(st, 0).Format(layout)
	case "H":
		layout = "2006010215"
		return "H-"+time.Unix(st, 0).Format(layout)
	case "I":
		layout = "200601021504"
		return "I-"+time.Unix(st, 0).Format(layout)
	}
	panic("unknown type of layout")
}