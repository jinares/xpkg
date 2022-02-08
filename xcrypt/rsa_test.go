package xcrypt

import (
	"encoding/hex"
	"fmt"
	"github.com/jinares/xpkg/xtools"
	"strconv"
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
	pri := "MIIEpAIBAAKCAQEAu9psw+Ddk++YO46uOpefrRjpCtm8Jaz4nerkxbexg6jKLINY0FS5AMsBsjNJ7NoCIkO82y+t0hqSG+A7+iZ+8ypuVgk5NNdHXtbpSV/IsGPybNt/AQxB/JHSrHOl7DdIZe6HlokQtFXiRlqZZQHcot4Q7iEMeXMzjnTXXpTcnWQsJuYOXJg+fW1CMtKyFEWvNowN90NJp6F3COKQ6BY7G+PA8Rf+goIEG1s4Fj8lN51huz8BczuLZDThN31lJBffrk/76jqhWfy6g8c4J++wvj0ie3GGhlpXq0vP/RYFwUfuIkzK0cUIGl/GyzmeqEK58mFhVZ9+Ltl3byVecLbpCwIDAQABAoIBABvP+AB8P3ISrCBNPKByKHLxWDeCCygs1tZSm7JmQwLyL5Stka2ujYF/jO28SDqoAF7gVorBIWQVCQLVjSu5HVyVMKbOUMjxAqEZAosZglYTRiY1zTTzgqWJVX6w1KUHflGRIkP4YR/CA3FVIpbg0rDOXgVToUwt8TkrNWPMeGlY1QcsyVEfIaoh8KtlB8omAmLzoy6PiqAGW+Lpsxcp9gsVjoQzM9LCARXGoVICbFQcIlo1j0oMipKRhru3tpv3H6UKADj0Tc9pI82SqtHHCxpSM7rlhPkcB00giHpJln3Lu1t9W4BWcdqYgAqm39bTAxbgVdj+8sJ5pmlgWO2Cx2kCgYEAwrEgXuMVml5Hcu16U8Vpjx3y46+1JHyFNH6MWnVd8914+3IY9UzxBTG+pqte3XhFD28J7D8zw3/LTlYrcNHBUFBeUHyiIZZ5c12Zlj7H1FXVBQpxvufs5PetgmrI2tuO9EyZxHegZKUZO6Wm89/Lym//CxV9MSIPVyHWtMyMghUCgYEA9wIBgq2AI+q1IWupfjo7svK8ugnZ8LqwuQWcoWJp1ZvF1O/YpbD+Tp+Istrhlyyj2HL+WaBQKTfp8scwR8FJxjbTld+3nI5Fo+UhWo7Lp0m9rPct9FlNbxU99edOwaSg23YefES4l7diLT5ZbZuCyLp25aL46VyQDxKi80FXJp8CgYEApX1gMae/Ji87dnJsB6cHWjKv0l/5jqEVrRBgh0e6a972xm0uz9vuB2dIUm3avlBMC5lsCteSTXxkORs/4684LeeSs8GtIvXAGJMYSUDmJRQsdRNyqj6D/ACYCQJx5q64bepqzjiNKt+3eh8NscCqflICrc4/UzNLbNoDoj04th0CgYEA4W04Oa3ka2MR6a1bI3M79qXnrZW+DCAllsZTW0n4stUWaK54N0df4Bti43A1QAWihrD0BpHzdpqr7UDyhBoYHUj+MyLYbI2/asN7fC0kGcmHzKpNi1pQ/BcT1C9Exh9cGs4jJmCFYxkfBZGIhirN4imixxLEPh2W79qfUogIZJsCgYB8f5XNzBtJv8PeSE/GEcrRufE6WhaZ5ZAJHeDWHaqF7XsX3Qcr3LZyNUUZwj5i5o6LCwtRURyj/JOaFhyPzDgq9y5g4AUgKAng7qYnX0KIe4yeAhZQBiZI8OH/Q0yDEfvpUNI/t4jZvVibBF5/k/Y0oMebGcIJK13YNTZtdjzO0g=="
	pub := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAu9psw+Ddk++YO46uOpefrRjpCtm8Jaz4nerkxbexg6jKLINY0FS5AMsBsjNJ7NoCIkO82y+t0hqSG+A7+iZ+8ypuVgk5NNdHXtbpSV/IsGPybNt/AQxB/JHSrHOl7DdIZe6HlokQtFXiRlqZZQHcot4Q7iEMeXMzjnTXXpTcnWQsJuYOXJg+fW1CMtKyFEWvNowN90NJp6F3COKQ6BY7G+PA8Rf+goIEG1s4Fj8lN51huz8BczuLZDThN31lJBffrk/76jqhWfy6g8c4J++wvj0ie3GGhlpXq0vP/RYFwUfuIkzK0cUIGl/GyzmeqEK58mFhVZ9+Ltl3byVecLbpCwIDAQAB"
	//pub, pri, err := GenRsaKeyWithPKCS1(2048 * 1)
	data, err := RsaSignPKCS1v15WithSHA256(pri, []byte("9"))
	fmt.Println(err)
	//fmt.Println(string(data))
	err = RsaVerfySignPKCS1v15WithSHA256([]byte("9"), data, pub)
	fmt.Println(err)
}

func TestHashMod(t *testing.T) {
	fmt.Println()
	fmt.Println("mod:", HashMod("F01ABCD444444444", 10))
	fmt.Println(strconv.ParseInt("F01ABCD4", 16, 64))
	fmt.Println(strconv.ParseInt("44444444", 16, 64))
	fmt.Println(HashMod("501", 100), fmt.Sprintf("%03d", HashMod("501", 200)))
	fmt.Println(HashMod("500", 100), fmt.Sprintf("%03d", HashMod("500", 200)))
	fmt.Println(HashMod("599", 100), fmt.Sprintf("%03d", HashMod("599", 200)))
	fmt.Println(HashMod("533", 100), fmt.Sprintf("%03d", HashMod("533", 200)))
	fmt.Println(HashMod("533", 100), fmt.Sprintf("%03d", HashMod("533", 200)))
	for i := 0; i < 10; i++ {
		gid := xtools.GUID()
		fmt.Println(len(gid), gid, HashMod(gid, 10))
	}
	fmt.Println(xtools.HashID("159518013", 100))

}
func TestRsaVerfySignPKCS1v15WithSHA256(t *testing.T) {

	data := `KEY=3001046&SIGN=9d8998922a0a204920790ef807134ec9&activeid=3001046&bizno=vr1203&ext1=bizno=000001203&fNum=1&mAct=3001046&ver=v3&ext2=front_order_time%3D1602680131859%26page_from%3Dwx_h5_pay%26trans_id%3D85402&ext3={"version":"v3","user_id":"10002004","rmb":"0.01","currency":"0","monthly":"0","request_type":"2","mBizNo":"000001203","mName":"皮玩语音金币","mDes":"皮玩语音金币","mAid":"99","mFgUrl":"https://www.youyisia.cn/recharge","mPeerId":"paycenter|other","mBgUrl":"http://xluser-web-paycallback.pay.svc/paycallback/v3/active","ip":"172.20.18.209","discount_amt":"0","fare_amt":"0","fact_amt":"1"}&fordertime=2020-10-14 20:55:31&fpaytype=N2&fxlpayid=201014205531131963T61236e2cf65&fxlpaytime=2020-10-14 20:56:04&num=1&orderid=20201014988520553151223845&other=&productid=3001046&rmb=1&userid=10002004`

	sm := "5e25189784eaeaf5063ef1ec6031bb8bc69782edcec307e4b8100aecc7a6db2b1a5e91d39ee0d7e8fb38cbb31880f8ef4ab5842be150a6b7bd768f5bf8887b4eac3046d75b7d42d943686843d83add1fb31007fdfdcd93f05fd030cc82f6a450656507704967407347fd597dfa77b08de0d3d8d81a6697fc55b77683495dc681a2fd7c457c07d123279fa5016f22b680a1e6a1c5de6690af6d27687bbbf5eb0d8faea24bfa854744dd358e56c898ab01f7419da2b5decce2e795168b2b94b3df589da018b651a4924a36d21c96831cbbc9d6ab89a85db79be4ed97ea16bca75e3309f20b59b87aad0a78cff38757529b8a9328a2c51e42264905bac9f3aeaf52"

	puk := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwI/Hye9L94vM1kKhvW5eB/ExBfg2Y5wmwVYqjNvRYNdiobK++l7wOr/t/VeJrfW6ivAwjB6St+qE4nHeU95lFgAu+N+GivVQp8uOsw1jKb2mbPbRt65ml7EGYUHkrGe0EKLYm0Z6gxoNg23KE9cbRgbK2/9j7Sf95AISW+kfr4yWgslZdf3mB/A7/lnwYp9vrCCze1XrUMUHEPPhaGB8yBADdoirr8ePr8FLJZrh8QIyeRO00ldtGREFuJD3CpjbXP8haOIpURuXWDijnbvqUtAg5ajJV2UoQ5iiV3ZpZShM66/BTRlb++KxvnMX6g131/U88CT8vgMibnAC7TnfvQIDAQAB"
	sm1, _ := hex.DecodeString(sm)
	ee := RsaVerfySignPKCS1v15WithSHA256([]byte(data), sm1, puk)
	fmt.Println(ee)
}
