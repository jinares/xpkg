package xcrypt

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"strconv"
)

//MD5  对一个字符串进行MD5加密,不可解密
func MD5(s string) string {
	h := md5.New()
	h.Write([]byte(s)) //使用zhifeiya名字做散列值，设定后不要变
	return hex.EncodeToString(h.Sum(nil))
}

//Sha1 获取 SHA1 字符串
func Sha1(s string) string {
	t := sha1.New()
	t.Write([]byte(s))

	return hex.EncodeToString(t.Sum(nil))
}

func HashMod(str string, mod int64) (ret int64) {
	defer func() {
		if ret == 0 {
			ret = mod
		}
	}()
	val, err := strconv.ParseInt(str, 10, 64)
	if err == nil && val > 0 {
		return val % mod
	}
	vi, err := strconv.ParseInt(str, 16, 64)
	if err == nil && vi > 0 {
		return vi % mod
	}
	md5_str := MD5(str)
	iad1, _ := strconv.ParseInt(md5_str[0:8], 16, 64)
	iad2, _ := strconv.ParseInt(md5_str[6:16], 16, 64)

	iad3, _ := strconv.ParseInt(md5_str[16:24], 16, 64)
	iad4, _ := strconv.ParseInt(md5_str[24:], 16, 64)

	return (((iad1%mod+iad2)%mod+iad3)%mod + iad4) % mod

}
