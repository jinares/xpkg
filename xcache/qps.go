package xcache

import (
	"sync"
	"sync/atomic"
	"time"
)

type (
	qpsItem struct {
		Val int64
	}
	qps struct {
		sync.RWMutex
		data    map[int64]*qpsItem
		ttl     int64
		isclose bool
	}
)

func NewQqs() *qps {
	ob := &qps{
		data:    map[int64]*qpsItem{},
		ttl:     3600,
		isclose: false,
	}
	go func() {
		for {
			time.Sleep(300 * time.Second)
			if ob.isclose {
				break
			}

			dt := time.Now().Unix()
			dl := []int64{}
			ob.RLock()
			for item_expire, _ := range ob.data {
				if (item_expire + ob.ttl) < dt {
					dl = append(dl, item_expire)
				}
			}
			ob.RUnlock()

			ob.Lock()
			for _, item := range dl {
				delete(ob.data, item)
			}
			ob.Unlock()
		}
	}()
	return ob
}
func (c *qps) Close() error {
	c.isclose = true
	c.Lock()
	c.data = map[int64]*qpsItem{}
	c.Unlock()
	return nil
}
func (c *qps) Incr() error {
	c.Lock()
	defer c.Unlock()
	dt := time.Now().Unix()
	val, isok := c.data[dt]
	if isok {
		atomic.AddInt64(&val.Val, 1)
		return nil
	}
	c.data[dt] = &qpsItem{
		Val: 1,
	}
	return nil
}
func (c *qps) GetQPS(sec int64) int64 {
	if sec > c.ttl {
		sec = c.ttl
	}
	dt := time.Now().Unix()
	all := int64(0)
	c.RLock()
	defer c.RUnlock()
	for i := int64(0); i < sec; i++ {
		val, isok := c.data[dt-i]
		if isok {
			all = all + val.Val
		}
	}
	return all / sec
}
