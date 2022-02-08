package xtools

import (
	"fmt"
	"testing"
)

func TestXMLToJSON(t *testing.T) {
	data := `
<xml>
	<ret>200</ret>
</xml>
`
	var ret struct {
		Ret int64 `json:"ret" xml:"ret"`
	}
	err := XMLToJSON(data, &ret)
	fmt.Println(err)
	fmt.Println(JSONToStr(ret))
	dd := `
`
	fmt.Println(dd)
}
