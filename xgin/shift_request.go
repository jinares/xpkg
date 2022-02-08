package xgin

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/gin-gonic/gin"
	"github.com/jinares/xpkg/xtools"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	//User-Agent
	USER_AGENT string = "xpkg/v2"
)

//DefaultTransport http keepalive 30s var tp http.RoundTripper
var DefaultTransport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   1 * time.Second,
		KeepAlive: 30 * time.Second,

		DualStack: true,
	}).DialContext,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
	TLSHandshakeTimeout:   3 * time.Second,
}

func newTransportWithTimeout(timeout time.Duration) *http.Transport {
	return &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   1 * time.Second,
			KeepAlive: 30 * time.Second,
			Deadline:  time.Now().Add(timeout * time.Second),
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSHandshakeTimeout:   3 * time.Second,
	}
}

func ShiftRequest(c *gin.Context, api string, sec int, retry int, header xtools.ES) (result string, err error) {
	//ctx := GinContext(c)
	if strings.EqualFold(c.Request.Method, "get") {
		return get(c, api, GetPar(c), sec, retry, header)
	} else if strings.EqualFold(c.Request.Method, "post") {
		return mPost(c, api, sec, retry, header)
	}
	return "", xtools.XErr(codes.InvalidArgument, "InvalidArgument", true)
}
func mPost(c *gin.Context, api string, sec int, retry int, header xtools.ES) (string, error) {
	ctx := GinContext(c)
	ctype := strings.ToLower(c.ContentType())
	par := ReadBody(c)

	if strings.HasPrefix(ctype, string(APPLICATION_URLENCODED)) {

		return postUrlencode(ctx, xtools.MergeUrl(api, GetPar(c)), par, sec, retry, header)
	}
	return postJSON(ctx, xtools.MergeUrl(api, GetPar(c)), par, sec, retry, header)

}

//get  get
func get(c *gin.Context, apiurl string, data map[string]string, second int, retry int, header xtools.ES) (string, error) {
	ctx := GinContext(c)
	api, err := url.Parse(apiurl)
	if err != nil {
		return "", xtools.XErr(codes.InvalidArgument, err.Error())
	}

	span, _ := opentracing.StartSpanFromContext(ctx, api.Scheme+"://"+api.Host+api.Path)
	defer span.Finish()
	req := httplib.Get(apiurl)
	opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.GetRequest().Header),
	)

	if strings.HasPrefix(strings.ToLower(apiurl), "https") {
		req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	req.Header("User-Agent", USER_AGENT)
	for k, v := range header {
		req.Header(k, v)
	}
	req.SetTransport(DefaultTransport)
	req.SetTimeout(1*time.Second, time.Duration(second)*time.Second)
	req.Retries(retry)
	for key, val := range data {
		req.Param(key, val)
	}
	start := time.Now()
	ret, err := req.String()
	if err != nil {
		return "", xtools.MErr(
			err, codes.DeadlineExceeded,
			fmt.Sprintf("timeout:%d", time.Since(start).Milliseconds()),
		)
	}
	return ret, nil
}

//Post  Post urlencode
func postUrlencode(ctx context.Context, apiurl string, data []byte, second int, retry int, header xtools.ES) (string, error) {
	api, err := url.Parse(apiurl)
	if err != nil {
		return "", err
	}

	span, _ := opentracing.StartSpanFromContext(ctx, api.Scheme+"://"+api.Host+api.Path)

	defer span.Finish()

	req := httplib.Post(apiurl)
	opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.GetRequest().Header),
	)

	if strings.HasPrefix(strings.ToLower(apiurl), "https") {
		req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	//req.SetTransport(newTransportWithTimeout(timeout))

	req.Header("User-Agent", USER_AGENT)
	req.Header("Content-Type", string(APPLICATION_URLENCODED))
	for k, v := range header {
		req.Header(k, v)
	}
	req.SetTransport(DefaultTransport)
	req.SetTimeout(1*time.Second, time.Duration(second)*time.Second)
	req.Retries(retry)
	req.Body(data)

	start := time.Now()
	ret, err := req.String()
	if err != nil {
		return "", xtools.MErr(
			err, codes.DeadlineExceeded,
			fmt.Sprintf("timeout:%d", time.Since(start).Milliseconds()),
		)
	}
	return ret, nil
}

//postJSON  Post  json
func postJSON(ctx context.Context, apiurl string, data []byte, second int, retry int, header xtools.ES) (string, error) {
	api, err := url.Parse(apiurl)
	if err != nil {
		return "", err
	}

	span, _ := opentracing.StartSpanFromContext(ctx, api.Scheme+"://"+api.Host+api.Path)
	defer span.Finish()

	req := httplib.Post(apiurl)
	opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.GetRequest().Header),
	)

	if strings.HasPrefix(strings.ToLower(apiurl), "https") {
		req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	//req.SetTransport(newTransportWithTimeout(timeout))
	req.SetTransport(DefaultTransport)
	req.SetTimeout(1*time.Second, time.Duration(second)*time.Second)
	req.Retries(retry)
	req.Header("User-Agent", USER_AGENT)
	req.Header("Content-Type", string(APPLICATION_JSON))
	for k, v := range header {
		req.Header(k, v)
	}
	req.Body(data)
	start := time.Now()
	ret, err := req.String()
	if err != nil {
		return "", xtools.MErr(
			err, codes.DeadlineExceeded,
			fmt.Sprintf("timeout:%d", time.Since(start).Milliseconds()),
		)
	}
	return ret, nil
}

func ShiftSvcPath(c *gin.Context, svc string) (string, error) {
	data, err := url.Parse(svc)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s://%s%s", data.Scheme, data.Host, c.Request.URL.Path), nil
}
