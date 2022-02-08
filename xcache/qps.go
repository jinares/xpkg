package xcache

import (
	"sync"
	"sync/atomic"
	"time"
)

type (
	QPSItem struct {
		Val int64
	}
	QPS struct {
		sync.RWMutex
		data    map[int64]*QPSItem
		ttl     int64
		isclose bool
	}
)

var (
	qpsMangerCtrlMutex sync.Once
	qpsMangerMutex     sync.RWMutex
	qpsManger          = map[string]*QPS{}
)

func NewQqs(name string, ttl int64) *QPS {
	if ttl < 0 {
		ttl = 3600
	}
	qpsMangerMutex.RLock()
	ob, isok := qpsManger[name]
	qpsMangerMutex.RUnlock()
	if isok {
		return ob
	}
	qpsMangerMutex.Lock()

	ob = &QPS{
		data:    map[int64]*QPSItem{},
		ttl:     3600,
		isclose: false,
	}
	qpsManger[name] = ob
	qpsMangerMutex.Unlock()
	qpsMangerCtrlMutex.Do(func() {
		go func() {
			for {
				time.Sleep(300 * time.Second)
				dt := time.Now().Unix()
				qpsMangerMutex.Lock()
				for _, qpsm := range qpsManger {
					qpsm.RLock()
					dl := []int64{}
					for item_expire, _ := range qpsm.data {
						if (item_expire + qpsm.ttl) < dt {
							dl = append(dl, item_expire)
						}
					}
					qpsm.RUnlock()
					if len(dl) == 0 {
						continue
					}
					qpsm.Lock()
					for _, it := range dl {
						delete(qpsm.data, it)
					}
					qpsm.Unlock()
				}
				qpsMangerMutex.Unlock()
			}
		}()
	})

	return ob
}
func (c *QPS) Close() error {
	c.isclose = true
	c.Lock()
	c.data = map[int64]*QPSItem{}
	c.Unlock()
	return nil
}
func (c *QPS) Incr() error {
	c.Lock()
	defer c.Unlock()
	dt := time.Now().Unix()
	val, isok := c.data[dt]
	if isok {
		atomic.AddInt64(&val.Val, 1)
		return nil
	}
	c.data[dt] = &QPSItem{
		Val: 1,
	}
	return nil
}
func (c *QPS) GetQPS(sec int64) int64 {
	if sec > c.ttl {
		sec = c.ttl
	}
	dt := time.Now().Unix()
	var all int64 = 0
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
