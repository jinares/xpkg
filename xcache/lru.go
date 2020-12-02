package xcache

import (
	"container/list"
	"github.com/jinares/xpkg/xerr"
	"google.golang.org/grpc/codes"
	"sync"
	"sync/atomic"
	"time"
)

type (
	lru struct {
		sync.RWMutex
		hash map[string]*list.Element

		expire   map[string]int64
		data     *list.List
		one      sync.Once
		capacity int64
		size     int64

		load AutoLoadHandler
	}
	AutoLoadHandler func(key string) (interface{}, int64, error)
)

func NewLRU(size int64) *lru {
	lru := &lru{
		hash:     map[string]*list.Element{},
		data:     list.New(),
		capacity: size, size: 0,
		expire: map[string]int64{},
	}

	lru.one.Do(func() {
		go func() {
			for {
				dt := time.Now().Unix()
				dlk := []string{}

				lru.RLock()

				for key, expire := range lru.expire {
					if dt <= expire {
						continue
					}
					dlk = append(dlk, key)
				}
				lru.RUnlock()
				if len(dlk) == 0 {
					time.Sleep(30 * time.Second)
					continue
				}
				lru.Lock()
				for _, key := range dlk {
					lru.del(key)
				}
				lru.Unlock()
			}
		}()
	})
	return lru
}
func (h *lru) Close() {
	h.Lock()
	defer h.Unlock()
	h.hash = map[string]*list.Element{}
	//h.expire = map[string]int64{}
	h.data = list.New()
	h.size = 0
}
func (h *lru) Len() int64 {
	return h.size
}
func (h *lru) Del(key string) error {
	h.Lock()
	defer h.Unlock()
	return h.del(key)

}
func (c *lru) del(key string) error {
	val, isok := c.hash[key]
	if isok == false {
		return nil
	}
	delete(c.hash, key)
	delete(c.expire, key)
	c.data.Remove(val)
	atomic.AddInt64(&c.size, -1)
	return nil
}

func (h *lru) set(key string, val interface{}, ttl int64) error {

	iv := &ItemVal{
		Key:    key,
		Val:    val,
		Expire: 0,
	}
	if ttl > 0 {
		iv.Expire = time.Now().Unix() + ttl
		h.expire[iv.Key] = iv.Expire
	}
	ele := h.data.PushBack(iv)
	h.hash[key] = ele
	next := atomic.AddInt64(&h.size, 1)
	if next > h.capacity {
		fele := h.data.Front()
		if fele != nil {
			val, isok := fele.Value.(*ItemVal)
			if isok {
				if val.Expire > 0 {
					delete(h.expire, val.Key)
				}
				//
				delete(h.hash, val.Key)
			}
			h.data.Remove(fele)
			atomic.AddInt64(&h.size, -1)
		}
	}
	return nil
}
func (h *lru) Setx(key string, val interface{}, ttl int64) error {
	h.Lock()
	defer h.Unlock()
	return h.set(key, val, ttl)
}
func (h *lru) Set(key string, val interface{}) error {
	return h.Setx(key, val, 0)
}
func (c *lru) auto_load_data(key string) (interface{}, error) {
	if c.load == nil {
		return nil, xerr.XErr(codes.NotFound, "not found")
	}
	v, exripe, err := c.load(key)
	if err != nil {
		return nil, err
	}
	if exripe < 0 {
		exripe = 0
	}
	c.set(key, v, exripe)
	return v, nil
}
func (h *lru) Get(key string) (interface{}, error) {
	h.RLock()
	ele, isok := h.hash[key]

	if isok {
		h.data.MoveToBack(ele)
	}
	h.RUnlock()
	if isok == false {
		h.Lock()
		defer h.Unlock()
		return h.auto_load_data(key)
	}
	if isok == false {
		return nil, xerr.XErr(codes.NotFound, "")
	}

	switch vv := ele.Value.(type) {
	case *ItemVal:
		if vv.Expire > 0 && vv.Expire < time.Now().Unix() {
			h.Lock()
			defer h.Unlock()
			return h.auto_load_data(key)
		}
		return vv.Val, nil

	default:
		return nil, xerr.XErr(codes.NotFound, "not found")
	}
}
func (h *lru) SetLoad(f AutoLoadHandler) *lru {
	if f == nil {
		return h
	}
	h.load = f
	return h
}
