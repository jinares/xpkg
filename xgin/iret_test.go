package xgin

import (
	"fmt"
	"net/url"
	"testing"
)

func TestRet_GetData(t *testing.T) {

	data, err := url.Parse("http://172.30.30.89:8000/?cmd=vip_adjust_order&b_id=2")
	fmt.Println(data, err, data.RequestURI(), data.Hostname(), data.String())
	fmt.Println(fmt.Sprintf("%s://%s%s", data.Scheme, data.Host, data.Path))
	fmt.Println(fmt.Sprintf("%s://%s%s", data.Scheme, data.Host, data.Path), data.Host, data.Hostname())

}
