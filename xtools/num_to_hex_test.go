package xtools

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func TestBHex2Num(t *testing.T) {
	v, err := BHex2Num("9", 10)
	fmt.Println(v, err)
	dt := time.Now()
	m := 5
	fmt.Println(dt.Format("2006"), 2020/m, 2021/m, 2025/m, 2026/m, 2030/m)
	fmt.Println(HashID("172503041", 100))
	dt, err = time.ParseInLocation("2006-01-02", "3000-01-02", time.Local)
	fmt.Println(err)
	fmt.Println(dt.Format(time.RFC3339))
	var a int64 = 127
	a = a << 40
	fmt.Println(a)
	fmt.Println(time.Now().Unix(), a+time.Now().Unix())
	fmt.Println(time.Unix(a+time.Now().Unix(), 0).Format(time.RFC3339))
	fmt.Println("maxint32", math.MaxInt32, math.MaxInt8)
	fmt.Println(time.Unix(int64(math.MaxUint32), 0).Format(time.RFC3339))
	fmt.Println(time.Unix(math.MaxInt64, 0).Format(time.RFC3339))

}
