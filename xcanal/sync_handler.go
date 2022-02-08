package xcanal

import (
	"regexp"
)

type (
	SyncHandler struct {
		DB      string `json:"db"`
		Table   string `json:"table"`
		Handler RowDataHandler
	}
	SyncHandlers []SyncHandler
)

func (h *SyncHandlers) Run(data *RowData) error {
	for _, item := range *h {
		err := item.Run(data)
		if err != nil {
			return err
		}
	}
	return nil
}
func (h *SyncHandler) check(data *RowData) bool {

	ret, err := regexp.MatchString(h.DB, data.GetDB())
	if err != nil || ret == false {
		return false
	}
	ret, err = regexp.MatchString(h.Table, data.GetTable())
	if err != nil || ret == false {
		return false
	}
	return true
}

func (h *SyncHandler) Run(row *RowData) error {
	if h.check(row) == false {
		return nil
	}
	if h.Handler == nil {
		return nil
	}
	return h.Handler(row)
}
func NewWriteHandler(db, table string, handler RowDataHandler) *SyncHandler {
	return &SyncHandler{
		DB:      db,
		Table:   table,
		Handler: handler,
	}
}
