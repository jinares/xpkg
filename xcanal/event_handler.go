package xcanal

import (
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
)

type (
	DefaultEventHandler struct {
		canal.DummyEventHandler

		fn         RowDataHandler
		c          *canal.Canal
		binlogName string
	}
)

func run(fn RowDataHandler, data *RowData) error {
	if fn == nil {
		return nil
	}
	return fn(data)
}

func (c *DefaultEventHandler) OnDDL(pos mysql.Position, queryEvent *replication.QueryEvent) error {
	if c.binlogName != pos.Name {
		c.binlogName = pos.Name
	}
	return run(c.fn, &RowData{
		Type: RowTypePos,
		Position: &Position{
			Name: pos.Name,
			Pos:  uint64(pos.Pos),
			GTID: "",
		},
	})
}
func (c *DefaultEventHandler) OnXID(pos mysql.Position) error {
	if c.binlogName != pos.Name {
		c.binlogName = pos.Name
	}
	return run(c.fn, &RowData{
		Type: RowTypePos,
		Position: &Position{
			Name: pos.Name,
			Pos:  uint64(pos.Pos),
			GTID: "",
		},
	})
}
func (c *DefaultEventHandler) OnGTID(mysql.GTIDSet) error {
	return nil
}
func (c *DefaultEventHandler) OnRotate(e *replication.RotateEvent) error {
	p := &Position{
		Name: string(e.NextLogName),
		Pos:  e.Position,
		GTID: "",
	}
	if c.binlogName != p.Name {
		c.binlogName = p.Name
	}

	return run(c.fn, &RowData{
		Type:     RowTypePos,
		Position: p,
		Data:     nil,
		Event:    nil,
	})
}
func (c *DefaultEventHandler) OnPosSynced(pos mysql.Position, gtid mysql.GTIDSet, isok bool) error {
	p := &Position{
		Name: pos.Name,
		Pos:  uint64(pos.Pos),
		GTID: "",
	}
	if gtid != nil {
		p.GTID = gtid.String()
	}
	if c.binlogName != pos.Name {
		c.binlogName = pos.Name
	}

	return run(c.fn, &RowData{
		Type:     RowTypePos,
		Position: p,
		Data:     nil,
		Event:    nil,
	})
}
func (c *DefaultEventHandler) OnRow(e *canal.RowsEvent) error {
	//fmt.Println(xtools.JSONToStr(e.Rows))
	data := []map[string]interface{}{}
	for _, val := range e.Rows {
		it := map[string]interface{}{}
		for i, v := range val {
			it[e.Table.Columns[i].Name] = v
		}
		data = append(data, it)
	}

	rd := &RowData{
		Type:  RowTypeData,
		Data:  data,
		Event: e,
	}
	if e.Header != nil && c.binlogName != "" {
		//fmt.Println(data)
		rd.Position = &Position{
			Name: c.binlogName, Pos: uint64(e.Header.LogPos),
		}
	}

	return run(c.fn, rd)
}

func (c *DefaultEventHandler) String() string {
	return "DefaultEventHandler"
}
