package xtools

import (
	"fmt"
	"testing"
)

func TestToYAML(t *testing.T) {
	data, err := ToYAML(map[string]interface{}{
		"data": GUID(),
		"mm": map[string]interface{}{
			"a": "a",
			"b": "b",
		},
	})
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(data, err)
}
