package xtools

import (
	"bytes"
	"encoding/json"
	"strings"
)

//JSONToStr JsonToStr
func JSONToStr(dstr interface{}) string {
	switch vv := dstr.(type) {
	case string:
		return vv
	}
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.Encode(dstr)

	return strings.TrimSpace(string(buffer.Bytes()))

}
func JSONMarshal(data interface{}) (string, error) {
	ret, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(ret), nil
}

//ToJson ToJson
func ToJSON(dstr string, data interface{}) error {
	err := json.Unmarshal([]byte(dstr), &data)
	return err
}
func DeepCopy(src interface{}, out interface{}) error {
	return ToJSON(JSONToStr(src), out)
}
