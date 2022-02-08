package xcanal

import (
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/jinares/xpkg/xtools"
	"strings"
)

type (
	//RowData RowData
	RowData struct {
		Type     RowType                  `json:"type"`
		Position *Position                `json:"position"`
		Data     []map[string]interface{} `json:"data"`
		Event    *canal.RowsEvent         `json:"event"`
	}
	RowType        int
	RowDataHandler func(row *RowData) (err error)

	//Position Position
	Position struct {
		Name string `yaml:"Name" json:"Name" bson:"Name"`
		Pos  uint64 `yaml:"Pos" json:"Pos" bson:"Pos"`
		GTID string `yaml:"GTID" json:"GTID" bson:"GTID"`
	}
)

const (
	RowTypeData      RowType = 0
	RowTypePos       RowType = 1
	RowTypeDump      RowType = 2
	RowTypeReSetGTID RowType = 3
)

func (c RowData) GetDB() string {
	if c.Type != RowTypeData {
		return ""
	}
	if c.Event == nil || c.Event.Table == nil {
		return ""
	}
	return c.Event.Table.Schema
}

func (c RowData) GetTable() string {
	if c.Type != RowTypeData {
		return ""
	}
	if c.Event == nil || c.Event.Table == nil {
		return ""
	}
	return c.Event.Table.Name
}
func (c *RowData) Count() int {
	return len(c.Data)
}
func (c *RowData) GetData(out interface{}) error {
	return xtools.DeepCopy(c.Data, out)
}
func (c *RowData) GetRet() []map[string]interface{} {
	if c.Type == RowTypeData {
		if c.UpdateAction() {
			return []map[string]interface{}{
				c.Data[0],
			}
		}
		return c.Data
	}
	return []map[string]interface{}{}
}

func (c *RowData) UpdateAction() bool {
	if strings.EqualFold(c.Event.Action, canal.UpdateAction) {
		if len(c.Data) == 2 {
			return true
		}
	}
	return false
}
func (c *RowData) DeleteAction() bool {
	if strings.EqualFold(c.Event.Action, canal.DeleteAction) {
		return true
	}
	return false
}
func (c *RowData) InsertAction() bool {
	if strings.EqualFold(c.Event.Action, canal.InsertAction) {
		return true
	}
	return false
}
func GetPK(pk []string, data map[string]interface{}) string {
	ret := []string{}
	for _, p := range pk {
		if val, isok := data[p]; isok {
			ret = append(ret, xtools.JSONToStr(val))
		}
	}
	return strings.Join(ret, ",")
}
