package xcache

import (
	"fmt"
	"github.com/allegro/bigcache/v3"
	"github.com/bluele/gcache"
	"github.com/jinares/xpkg/xtools"
	"google.golang.org/grpc/codes"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestLRU_Set(t *testing.T) {
	/*
		Shards：切片数，必须为2的整数幂；
		LifeWindow：数据可被清除(驱逐)的时间；
		CleanWindow：扫描清除的时间间隔；
		MaxEntriesInWindow：初始的数据规模；
		MaxEntrySize：初始的数据字节大小；
		Verbose：是否输出运行信息；
		HardMaxCacheSize：硬性的最大规模；
		OnRemove：清除数据的回调；
		OnRemoveWithReason：包括清除原因的清除回调。
	*/
	cfg := bigcache.DefaultConfig(1 * time.Second)
	cfg.HardMaxCacheSize = 3
	cfg.CleanWindow = 0
	cache, _ := bigcache.NewBigCache(cfg)

	//cache.Set("my-unique-key", []byte("value"))
	for i := 0; i <= 200000; i++ {

		cache.Set(xtools.GUID(), []byte(xtools.GUID()))
	}
	cache.Set("my-unique-key", []byte("value"))
	fmt.Println("status:", cache.Capacity(), cache.Len())
	time.Sleep(5 * time.Second)

	cache.Set("my-unique-key2", []byte("value"))
	entry, err := cache.Get("my-unique-key")
	fmt.Println(string(entry), err, cache.Capacity(), cache.Len())

	var ii int32 = 0
	ne := atomic.AddInt32(&ii, 1)
	atomic.AddInt32(&ii, 1)
	fmt.Println(ii, ne, len("20211015115048629081503495bizmobilegame"))

}
func mKey(userid string, item_id int64) string {
	return fmt.Sprintf("%s#%d", userid, item_id)
}
func decodeKey(key string) (userid string, itemid int64, err error) {
	if key == "" {
		return "", 0, xtools.XErr(codes.InvalidArgument, "")
	}
	sl := strings.Split(key, "#")
	if len(sl) != 2 {
		return "", 0, xtools.XErr(codes.InvalidArgument, "")
	}
	ival, isok := xtools.IntVal(sl[1])
	if isok == false {

		return "", 0, xtools.XErr(codes.InvalidArgument, "")
	}
	return sl[0], ival, nil
}

type (
	parKey struct {
		UserID string
		ItemID int64
	}
)

func BenchmarkLRU_Set33(b *testing.B) {
	b.StopTimer()
	cfg := bigcache.DefaultConfig(10 * time.Minute)
	cfg.HardMaxCacheSize = 128
	cache, _ := bigcache.NewBigCache(cfg)
	for i := 0; i < 150000; i++ {
		cache.Set(xtools.GUID(), []byte(xtools.GUID()))
	}
	//cache.Set("my-unique-key", []byte("value"))

	//entry, _ := cache.Get("my-unique-key")
	//fmt.Println(entry)
	b.SetParallelism(100)
	b.StartTimer()
	for idx := 0; idx < b.N; idx++ {

		cache.Get(xtools.ToStr(idx))
	}
	b.ReportAllocs()
}
func BenchmarkLRU_Set(b *testing.B) {
	b.StopTimer()
	lru := NewLRU(300000)
	b.SetParallelism(100)
	b.StartTimer()
	for idx := 0; idx < b.N; idx++ {
		//lru.Set(xtools.GUID(), xtools.GUID())
		lru.Get(xtools.GUID())
	}
	b.ReportAllocs()
}
func BenchmarkGc_Set(b *testing.B) {
	b.StopTimer()
	gc := gcache.New(300000).LRU().Build()
	b.SetParallelism(100)
	b.StartTimer()
	for idx := 0; idx < b.N; idx++ {
		//gc.Set(xtools.GUID(), xtools.GUID())
		gc.Get(xtools.ToStr(idx))
	}
	b.ReportAllocs()

}
func BenchmarkLRU_Del(b *testing.B) {
	b.StopTimer()
	lru := NewLRU(20000)
	for i := 0; i < 20000; i++ {

		lru.Set(xtools.ToStr(i), xtools.GUID())
	}
	b.StartTimer()
	for idx := 0; idx < b.N; idx++ {
		lru.Del(xtools.ToStr(idx))
	}
	b.ReportAllocs()
}

func BenchmarkGc_Del(b *testing.B) {
	b.StopTimer()
	gc := gcache.New(20000).LRU().Build()
	for i := 0; i < 20000; i++ {
		gc.Set(xtools.GUID(), xtools.GUID())
	}
	b.StartTimer()
	for idx := 0; idx < b.N; idx++ {
		gc.Remove(xtools.ToStr(idx))
	}
	b.ReportAllocs()
}

func BenchmarkLRU_Len(b *testing.B) {
	b.StopTimer()
	lru := NewLRU(20000)
	for i := 0; i < 20000; i++ {

		lru.Set(xtools.ToStr(i), xtools.GUID())
	}
	b.StartTimer()
	for idx := 0; idx < b.N; idx++ {
		lru.Len()
	}
	b.ReportAllocs()
}

func BenchmarkGc_Len(b *testing.B) {
	b.StopTimer()
	gc := gcache.New(20000).LRU().Build()
	for i := 0; i < 20000; i++ {
		gc.Set(xtools.GUID(), xtools.GUID())
	}
	b.StartTimer()
	for idx := 0; idx < b.N; idx++ {
		gc.Len(false)
	}
	b.ReportAllocs()
}
func BenchmarkLRU_Get(b *testing.B) {
	b.StopTimer()
	lru := NewLRU(20000)
	for i := 0; i < 20000; i++ {

		lru.Set(xtools.ToStr(i), xtools.GUID())
	}
	b.StartTimer()
	for idx := 0; idx < b.N; idx++ {
		lru.Get(xtools.ToStr(idx))
	}
	b.ReportAllocs()
}

func BenchmarkGc_Get(b *testing.B) {
	b.StopTimer()
	gc := gcache.New(20000).LRU().Build()
	for i := 0; i < 20000; i++ {
		gc.Set(xtools.GUID(), xtools.GUID())
	}
	b.SetParallelism(200)
	b.StartTimer()
	for idx := 0; idx < b.N; idx++ {
		gc.Get(xtools.ToStr(idx))
		gc.Set(xtools.ToStr(idx), xtools.ToStr(idx))
	}
	b.ReportAllocs()
}
func BenchmarkLRU_SetX(b *testing.B) {
	b.StopTimer()
	b.SetParallelism(100)
	lru := NewLRU(300000)

	b.StartTimer()
	for idx := 0; idx < b.N; idx++ {
		lru.Setx(xtools.GUID(), xtools.GUID(), 3600)
	}
	b.ReportAllocs()
}

func BenchmarkGc_Setx(b *testing.B) {
	b.StopTimer()
	gc := gcache.New(300000).LRU().Build()
	b.SetParallelism(100)
	b.StartTimer()
	for idx := 0; idx < b.N; idx++ {
		gc.SetWithExpire(xtools.GUID(), xtools.GUID(), 3600)
	}
	b.ReportAllocs()

}
