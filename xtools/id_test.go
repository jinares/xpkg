package xtools

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestGUID(t *testing.T) {
	fmt.Println(GUID(), len(GUID()))
	b := xmd5(GUID())

	uuid := fmt.Sprintf("%s-%s-%s-%s-%s", b[0:8], b[8:12], b[12:16], b[16:20], b[20:])
	fmt.Println(strings.ToLower(uuid))
	fmt.Println(time.Now().Format(time.RFC3339))
	ctx := context.WithValue(context.Background(), "x-auth", GUID())
	//ctx = context.WithValue(ctx, "key", "123")

	fmt.Println(ctx)
	fmt.Println(ctx.Value("key"))
	fmt.Println(ctx.Value("x-auth"))
	fmt.Println(5%10, 5%5)
	fmt.Println(2&31, 5&31)

}

func TestHashID(t *testing.T) {
	dt := time.Now()
	m := 100
	id := GUID()
	k := xmd5(id)
	hash := int64(0)
	fmt.Println(k)
	for v := range k {
		va := int64(k[v])
		fh := strconv.FormatInt(va, 16)
		fo, _ := strconv.ParseInt(fh, 10, 64)
		hash += fo
	}
	fmt.Println(hash)
	hash = (hash * 1) % int64(m)
	fmt.Println(hash + 1)
	fmt.Println(HashID("793055748", 100))
	time.Sleep(100 * time.Millisecond)
	s := time.Since(dt)
	fmt.Println(s.Milliseconds())
}
