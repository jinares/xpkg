package xtools

import (
	"fmt"
	"time"

	//"net/url"
	"gopkg.in/yaml.v2"
	"testing"
)

type (
	tes struct {
		PayType    string `json:"PayType" yaml:"PayType"`
		AppID      string `json:"AppID" yaml:"AppID"`
		Account    string `json:"Account" yaml:"Account"`
		CreateTime string `json:"CreateTime" yaml:"CreateTime"`
		Num        int    `json:"Num" yaml:"Num"`
	}
)

func TestUrlEncode(t *testing.T) {
	var tet struct {
		Name string   `json:"name"  yaml:"Name"`
		List []tes    `json:"List" yaml:"List"`
		Tes  tes      `json:"tes" yaml:"Tes"`
		Tags []string `json:"Tags" yaml:"Tags,flow"`
	}
	tet.Name = "fdsaf"
	dt := time.Now().Format("2006-01-02 15:04:05")
	tet.List = []tes{
		{PayType: "AM", AppID: "1234577890", Account: GUID(), Num: 1, CreateTime: dt},
		{PayType: "AM", AppID: GUID(), Account: GUID(), Num: 2, CreateTime: dt},
		{PayType: "AM", AppID: GUID(), Account: GUID(), Num: 3, CreateTime: dt},
	}
	tet.Tags = []string{"aaa", "bbb", "ccc"}
	tet.Tes = tes{PayType: "AM", AppID: GUID(), Account: GUID(), Num: 3, CreateTime: dt}
	ret, _ := yaml.Marshal(tet)
	fmt.Println(len(string(ret)), len(JSONToStr(tet)))
	fmt.Println(string(ret))

}
