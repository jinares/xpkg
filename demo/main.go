package main

import (
	"fmt"
	"github.com/jinares/xpkg/xcrypt"
	"github.com/jinares/xpkg/xerr"
	"github.com/jinares/xpkg/xetcd"
	"github.com/jinares/xpkg/xlog"
	"github.com/jinares/xpkg/xtools"
	"google.golang.org/grpc/codes"
	"time"
)

func main() {
	fmt.Println(xcrypt.Base64Encode("1234"))
	xetcd.NewCli(nil, "/root")
	err := xerr.XErr(codes.InvalidArgument, "err", true)
	fmt.Println(err.String(), xtools.Caller(0, true))
	xlog.Info("xlog")
	data := "12345"
	fmt.Println(data[2:4])
	time.n
	fmt.Println(time.Now().Format(time.RFC3339))
}
func longestPalindrome(s string) string {
	l := len(s)
	if l <= 1 {
		return s
	}
	for i := 1; i < l; i++ {

	}

	return s

}
