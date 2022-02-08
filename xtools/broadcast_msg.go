package xtools

import (
	"container/list"
	"fmt"
	"google.golang.org/grpc/codes"
	"sync"
)

type (
	Broadcast interface {
		Close() error
		Push(data string) error
		Subscriber() string
	}

	BroadcastItem struct {
		data    chan string
		mgr     *BroadcastManager
		ele     *list.Element
		isclose bool
	}
	BroadcastManager struct {
		sync.RWMutex
		data    *list.List
		itemCap int
	}
)

func NewBroadcastManager(cap int) *BroadcastManager {
	if cap <= 0 {
		cap = 1
	}
	return &BroadcastManager{data: list.New(), itemCap: cap}
}
func (c *BroadcastItem) Push(data string) error {
	if c.isclose {
		fmt.Println("isclose.....")
		return nil
	}
	select {
	case c.data <- data:

		return nil

	default:
		return XErr(codes.ResourceExhausted, "push-timeout")

	}
}

func (c *BroadcastItem) Subscriber() string {
	if c.isclose {
		return ""
	}
	return <-c.data
}
func (c *BroadcastItem) Close() error {
	if c.ele == nil {
		return nil
	}
	if c.isclose {
		return nil
	}
	c.mgr.Lock()
	defer c.mgr.Unlock()
	c.isclose = true
	c.mgr.data.Remove(c.ele)

	return nil
}

func (c *BroadcastManager) Push(data string) error {
	//c.RLock()
	//defer c.RUnlock()
	for e := c.data.Front(); e != nil; e = e.Next() {
		switch br := e.Value.(type) {
		case Broadcast:
			err := br.Push(data)
			if err != nil {
				return err
			}
			return nil
		default:
			return nil
		}
	}
	return nil
}
func (c *BroadcastManager) Len() int {
	return c.data.Len()
}

func (c *BroadcastManager) Get() (Broadcast, error) {
	c.Lock()
	defer c.Unlock()
	m := &BroadcastItem{
		data: make(chan string, c.itemCap),
		mgr:  c,
	}
	el := c.data.PushBack(m)
	m.ele = el
	return m, nil

}
