package xcache

import (
	"fmt"
	"testing"
	"time"
)

func TestNewQqs(t *testing.T) {
	ob := NewQqs()
	for i := 0; i < 2000; i++ {
		ob.Incr()
	}
	time.Sleep(1 * time.Second)
	for i := 0; i < 1000; i++ {
		ob.Incr()
	}
	time.Sleep(1 * time.Second)
	for i := 0; i < 1000; i++ {
		ob.Incr()
	}
	time.Sleep(1 * time.Second)
	for i := 0; i < 1000; i++ {
		ob.Incr()
	}
	fmt.Println(ob.GetQPS(60), ob.GetQPS(3))
}
func BenchmarkNewQqs(b *testing.B) {
	b.StopTimer()
	ob := NewQqs()
	for i := 0; i < 1000; i++ {
		ob.Incr()
	}
	b.StartTimer()
	for idx := 0; idx < b.N; idx++ {
		ob.Incr()
	}

	b.ReportAllocs()

	//fmt.Println("total:", ob.GetQPS(60), ob.GetQPS(3))
}
