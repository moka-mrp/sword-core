package sign

const (
	HMACSHA256                 = "HMAC-SHA256"
	MD5                        = "MD5"
	SHA1                       = "SHA1"
	Sn                       = "sn"  //必传参数
	Timestamp                  = "timestamp" //必传参数
	Nonce                      = "nonce" //必传参数
	SignatureExpired=600 //签名有效期暂时为600s
	JwUrl             ="url"
)


