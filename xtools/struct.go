package xtools

import (
	"fmt"
	"reflect"
	"strings"
)

type (
	FieldTag struct {
		Value string   `json:"value"`
		Ext   []string `json:"ext"`
	}
)

func NewTag(tag string) FieldTag {
	tmp := strings.Split(tag, ",")
	if len(tmp) < 2 {
		return FieldTag{
			Value: tag,
			Ext:   []string{},
		}
	}
	return FieldTag{
		Value: tmp[0],
		Ext:   tmp[1:],
	}
}
func (h *FieldTag) Encode() string {
	if len(h.Ext) < 1 {
		return h.Value
	}
	return fmt.Sprintf("%s,%s", h.Value, strings.Join(h.Ext, ","))
}
func (h *FieldTag) HasExt(ext ...string) bool {
	tmp := map[string]string{}
	for _, val := range ext {
		tmp[val] = ""
	}
	for _, val := range h.Ext {
		if _, isok := tmp[val]; isok {
			return true
		}
	}
	return false
}

//StructToStringMap　struct => map[string]string
func StructToStringMap(data interface{}) ES {

	ret := ES{}
	ss := reflect.ValueOf(data)
	switch reflect.TypeOf(data).Kind() {
	case reflect.Map:
		//fmt.Println(data)
		mval := ss.MapRange()
		for {
			if mval.Next() == false {
				break
			}
			ret[JSONToStr(mval.Key().Interface())] = JSONToStr(mval.Value().Interface())
		}
		return ret
	case reflect.Slice:

		for i := 0; i < ss.Len(); i++ {
			ssval := ss.Index(i)
			if ssval.CanInterface() {
				ret[ToStr(i)] = JSONToStr(ssval.Interface())
			}

		}
		//fmt.Println("slice", data, ss.Index(0), ss.Len())
		return ret
	case reflect.Struct:

	case reflect.Ptr:
		ss = ss.Elem()
	default:
		return ES{}

	}
	tag := `json`
	//ss := reflect.ValueOf(data).Elem()
	tf := NewTag(tag)
	for i := 0; i < ss.NumField(); i++ {
		tmp := ss.Type().Field(i).Tag.Get(tf.Value)
		stf := NewTag(tmp)
		if len(tmp) <= 0 {
			continue
		}
		if len(tf.Ext) > 0 && tf.HasExt(stf.Ext...) == false {
			continue
		}
		name := ss.Type().Field(i).Name
		field := ss.FieldByName(name)
		if field.Kind() != reflect.String {
			if field.CanInterface() {
				ret[stf.Value] = JSONToStr(field.Interface())
			}

			continue
		}
		val := field.String()
		ret[stf.Value] = val

	}
	return ret
}
