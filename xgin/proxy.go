package xgin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type (
	FilterHandler func(ctx context.Context, opt FilterOption) bool
	FilterOption  struct {
		Path  string            `json:"path"`
		Query map[string]string `json:"query"`
		Body  []byte            `json:"body"`
	}
	ServerInfo struct {
		ServerdAddr          string `yaml:"ServerdAddr"`
		ServerTimeout        int    `yaml:"ServerTimeout"`
		KeepaliveIdleTimeout int    `yaml:"KeepaliveIdleTimeout"`
	}
	GinHandler func(xg *gin.Engine)
)

func RunEngine(sv ServerInfo, fn GinHandler) error {
	xg := gin.New()
	fn(xg)
	return RunRouter(xg, sv)
}
func RunRouter(engine *gin.Engine, sv ServerInfo) (err error) {

	s := &http.Server{
		Addr:        sv.ServerdAddr,
		Handler:     engine,
		ReadTimeout: time.Duration(sv.ServerTimeout) * time.Millisecond,
		IdleTimeout: time.Duration(sv.KeepaliveIdleTimeout) * time.Millisecond,
	}
	if sv.KeepaliveIdleTimeout <= 0 {
		s.SetKeepAlivesEnabled(false)
	}

	err = s.ListenAndServe()

	if err != nil {
		panic(err)
	}
	return
}
func (c *FilterOption) GetVal(key string) (string, bool) {
	val, isok := c.Query[key]
	return val, isok
}
func (c *FilterOption) GetBodyJSON(out interface{}) error {
	return json.Unmarshal(c.Body, out)
}
func ProxyURL(c *gin.Context, targeturl string, fn FilterHandler) bool {
	opt := FilterOption{
		Path:  c.Request.URL.Path,
		Query: GetPar(c),
		Body:  ReadBody(c),
	}
	if fn(GinContext(c), opt) == false {
		return false
	}
	data, err := url.Parse(targeturl)
	if err != nil {
		return false
	}
	data, _ = url.Parse(fmt.Sprintf("%s://%s%s", data.Scheme, data.Host, data.Path))
	if data.Path != "" {
		c.Request.URL.Path = data.Path
		data.Path = ""
	}
	//brw := newbResponseWriter(c.Writer)
	//c.Writer = brw

	p := httputil.NewSingleHostReverseProxy(data)
	p.ServeHTTP(c.Writer, c.Request)

	//brw.Flush()

	c.Abort()
	return true
}

//
//var DefaultTransport = &http.Transport{
//	Proxy: http.ProxyFromEnvironment,
//	DialContext: (&net.Dialer{
//		Timeout:   30 * time.Second,
//		KeepAlive: 30 * time.Second,
//		DualStack: true,
//	}).DialContext,
//	ForceAttemptHTTP2:     true,
//	MaxIdleConns:          100,
//	IdleConnTimeout:       90 * time.Second,
//	TLSHandshakeTimeout:   10 * time.Second,
//	ExpectContinueTimeout: 1 * time.Second,
//}
var ProxyTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: -1,
		DualStack: true,
	}).DialContext,
	ForceAttemptHTTP2:     true,
	DisableKeepAlives:     true,
	MaxIdleConns:          1,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}
var (
	SingleHostReverseProxy      = map[string]*httputil.ReverseProxy{}
	mutexSingleHostReverseProxy sync.RWMutex
)

func NewSingleHostReverseProxy(data *url.URL) *httputil.ReverseProxy {
	svc := fmt.Sprintf("%s://%s", data.Scheme, data.Host)
	mutexSingleHostReverseProxy.RLock()
	rp, isok := SingleHostReverseProxy[svc]
	mutexSingleHostReverseProxy.RUnlock()
	if isok {
		return rp
	}
	mutexSingleHostReverseProxy.Lock()
	defer mutexSingleHostReverseProxy.Unlock()
	rp, isok = SingleHostReverseProxy[svc]
	if isok {
		return rp
	}
	srp := httputil.NewSingleHostReverseProxy(data)
	srp.Transport = ProxyTransport
	SingleHostReverseProxy[svc] = srp

	return srp
}
func ProxyURLV2(c *gin.Context, targeturl string, fn FilterHandler) (string, bool) {
	opt := FilterOption{
		Path:  c.Request.URL.Path,
		Query: GetPar(c),
		Body:  ReadBody(c),
	}
	if fn(GinContext(c), opt) == false {
		return "", false
	}
	data, err := url.Parse(targeturl)
	if err != nil {
		return "", false
	}
	data, _ = url.Parse(fmt.Sprintf("%s://%s%s", data.Scheme, data.Host, data.Path))
	if data.Path != "" {
		c.Request.URL.Path = data.Path
		data.Path = ""
	}
	fmt.Println(data, c.Request.URL.Path, c.Request.URL.Host, c.Request.Host)
	p := httputil.NewSingleHostReverseProxy(data)
	var result string = ""
	p.ModifyResponse = func(resp *http.Response) error {
		rc, ret, err := ReadBodyv2(resp.Body)
		if err != nil {
			return err
		}
		resp.Body = rc
		result = ret
		return nil
	}

	p.ServeHTTP(c.Writer, c.Request)
	//runtime.GC()
	c.Abort()

	//brw := newbResponseWriter(c.Writer)
	//c.Writer = brw
	//p.ServeHTTP(c.Writer, c.Request)
	//result := brw.GetData()
	//brw.Flush()
	//c.Abort()
	return string(result), true
}
func ReadBodyv2(rc io.ReadCloser) (io.ReadCloser, string, error) {
	body, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, "", err
	}
	rc.Close()
	//rc = ioutil.NopCloser(bytes.NewBuffer(body))
	return ioutil.NopCloser(bytes.NewBuffer(body)), string(body), nil
}
func ReadBody(c *gin.Context) []byte {
	body, _ := ioutil.ReadAll(c.Request.Body)
	c.Request.Body.Close()
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return body
}

func CopyBody(c *gin.Context) (io.ReadCloser, error) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}
	c.Request.Body.Close()
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return ioutil.NopCloser(bytes.NewBuffer(body)), nil
}
