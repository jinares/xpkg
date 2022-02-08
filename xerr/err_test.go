package xerr

import (
	"fmt"
	"github.com/jinares/xpkg/xtools"
	"google.golang.org/grpc/codes"
	"testing"
)

func TestXErr(t *testing.T) {
	xe := XErr(codes.InvalidArgument, "test", true)
	fmt.Println(xtools.JSONToStr(xe))
	fmt.Println(ErrString(xe))
	fmt.Println(ErrMsg(xe))
	fmt.Println(xe.Error())
}
