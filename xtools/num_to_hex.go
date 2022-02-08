package xtools

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math"
	"strings"
)

const (
	num2char = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func NumToBHex(num, n int64) (string, error) {
	if n != 16 && n != 36 && n != 62 && n != 10 && n != 8 {
		return "", status.Error(codes.InvalidArgument, "")
	}
	if num == 0 {
		return "0", nil
	}
	num_str := ""
	for num != 0 {
		yu := num % n
		num_str = string(num2char[yu]) + num_str
		num = num / n
	}
	return num_str, nil
}
func BHex2Num(str string, n int) (int64, error) {
	if n != 16 && n != 36 && n != 62 && n != 10 && n != 8 {
		return 0, status.Error(codes.InvalidArgument, "")
	}
	//str = strings.ToLower(str)
	v := 0.0
	length := len(str)
	for i := 0; i < length; i++ {
		s := string(str[i])
		index := strings.Index(num2char, s)
		v += float64(index) * math.Pow(float64(n), float64(length-1-i)) // 倒序
	}
	return int64(v), nil
}
