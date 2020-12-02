package xcache

import (
	"fmt"
	"github.com/bluele/gcache"
	"github.com/jinares/xpkg/xtools"
	"testing"
)

func TestLRU_Set(t *testing.T) {
	lru := NewLRU(3).SetLoad(func(key string) (interface{}, int64, error) {
		return xtools.GUID(), 3600, nil
	})
	fmt.Println("len:", lru.Len())
	fmt.Println(lru.Get("3"))
	fmt.Println(lru.Get("2"))
	fmt.Println(lru.Get("1"))

	lru.Set("4", "4")

	fmt.Println("len:", lru.Len())
	fmt.Println(lru.Get("1"))
	fmt.Println(lru.Get("2"))
	fmt.Println(lru.Get("3"))

	fmt.Println(lru.Get("4"))
	fmt.Println(329853488333 / 1024 / 1024 / 1024)

}
func BenchmarkLRU_Set(b *testing.B) {
	b.StopTimer()
	lru := NewLRU(3000)

	b.StartTimer()
	for idx := 0; idx < b.N; idx++ {
		lru.Set(xtools.GUID(), xtools.GUID())
	}
	b.ReportAllocs()
}
func BenchmarkGc_Set(b *testing.B) {
	b.StopTimer()
	gc := gcache.New(300000).LRU().Build()
	b.SetParallelism(100)
	b.StartTimer()
	for idx := 0; idx < b.N; idx++ {
		gc.Set(xtools.GUID(), xtools.GUID())
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
	b.StartTimer()
	for idx := 0; idx < b.N; idx++ {
		gc.Get(xtools.ToStr(idx))
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
