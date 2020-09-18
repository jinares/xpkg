package xcrypt

import (
	"fmt"
	"testing"
)

func TestGenRsaKeyWithPKCS1(t *testing.T) {
	pub, pri, err := GenRsaKeyWithPKCS1(2048 * 4)
	fmt.Println(len(pub), pub)
	fmt.Println(len(pri), pri)
	fmt.Println(err)
	data := "5a4456e8fc1df8d82d1439e5a6259b35108b16be72a1caaa02edf209ba013ccd38fb5df606cf4249ab316c3dbf39fd0773b863"
	sign_msg, err := RsaSignPKCS1v15WithSHA256(pri, []byte(data))
	fmt.Println(string(sign_msg))
	fmt.Println(len(HexEncodeStr(string(sign_msg))), HexEncodeStr(string(sign_msg)))
	err = RsaVerfySignPKCS1v15WithSHA256([]byte(data), sign_msg, pub)
	fmt.Println(err)
}

func TestCRC32(t *testing.T) {
	fmt.Println(CRC32("123456789"))
}
