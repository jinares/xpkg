package xgin

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net"
)

const (
	// noWritten
	noWritten = -1
)

// 重写 HTTP输出
type bResponseWriter struct {
	buffer *bytes.Buffer
	gin.ResponseWriter
	size    int
	isFlush bool
	status  int
}

// output http code
func (rw *bResponseWriter) WriteHeader(code int) {
	rw.status = code
}

// output http header
func (rw *bResponseWriter) WriteHeaderNow() {
	rw.ResponseWriter.WriteHeaderNow()
}

// output body byte
func (rw *bResponseWriter) Write(data []byte) (n int, err error) {
	n, err = rw.buffer.Write(data)
	rw.size += n
	return
}

// output body string
func (rw *bResponseWriter) WriteString(s string) (n int, err error) {
	n, err = io.WriteString(rw.buffer, s)
	rw.size += n
	return
}

// get http status
func (rw *bResponseWriter) Status() int {
	return rw.ResponseWriter.Status()
}

// rewrite http body
func (rw *bResponseWriter) ReWrite(data []byte) (int, error) {
	rw.buffer.Reset()
	rw.size = 0

	return rw.Write(data)
}

// get size
func (rw *bResponseWriter) Size() int {
	return rw.size
}

// get Written
func (rw *bResponseWriter) Written() bool {
	return rw.size != noWritten
}

// Hijack
func (rw *bResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return rw.ResponseWriter.Hijack()
}

// CloseNotify
func (rw *bResponseWriter) CloseNotify() <-chan bool {
	return rw.ResponseWriter.CloseNotify()
}

// Flush
func (rw *bResponseWriter) Flush() {
	if rw.isFlush {
		return
	}
	rw.ResponseWriter.WriteHeader(rw.status)
	if rw.buffer.Len() > 0 {
		data := rw.buffer.Bytes()
		_, err := rw.ResponseWriter.Write(data)
		if err != nil {
			//panic(err)
			fmt.Println("bResponseWriter:", err.Error())
		}
		rw.buffer.Reset()
	}
	rw.isFlush = true
}
func (rw *bResponseWriter) GetData() []byte {
	return rw.buffer.Bytes()
}

var _ gin.ResponseWriter = &bResponseWriter{}

// newbResponseWriter
func newbResponseWriter(rw gin.ResponseWriter) *bResponseWriter {
	bresp := &bResponseWriter{ResponseWriter: rw, buffer: &bytes.Buffer{}, status: 200}

	return bresp
}
