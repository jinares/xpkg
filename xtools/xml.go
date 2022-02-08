package xtools

import (
	"encoding/xml"
)

func XMLToJSON(str string, obj interface{}) error {
	return xml.Unmarshal([]byte(str), obj)
}
