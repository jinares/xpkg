package xtools

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/astaxie/beego/httplib"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/url"
	"strings"
	"time"
)

const (
	//User-Agent
	USER_AGENT string = "ares/v1"
)

//MergeUrl 合并URL参数
func MergeUrl(uri string, param map[string]string) string {

	ret, err := url.Parse(uri)
	if err != nil {
		if strings.Contains(uri, "?") {
			return uri + "&" + UrlEncode(param)
		}
		return uri + "?" + UrlEncode(param)
	}
	q := UrlDecode(ret.RawQuery)
	for key, val := range param {
		q[key] = val
	}
	ret.RawQuery = UrlEncode(q)
	return ret.String()
}

//Get  Get
func Get(ctx context.Context, apiurl string, data map[string]string, second int, retry int) (string, error) {
	api, err := url.Parse(apiurl)
	if err != nil {
		return "", status.Error(codes.InvalidArgument, err.Error())
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
	req.SetTimeout(3*time.Second, time.Duration(second)*time.Second)
	req.Retries(retry)
	for key, val := range data {
		req.Param(key, val)
	}
	start := time.Now()

	ret, err := req.String()
	if err != nil {
		return "", status.Error(codes.DeadlineExceeded, ToStr(time.Since(start).Seconds())+" msg:"+err.Error())
	}
	return ret, nil
}

//Post  Post urlencode
func Post(ctx context.Context, apiurl string, data map[string]string, second int, retry int) (string, error) {
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
	req.SetTimeout(3*time.Second, time.Duration(second)*time.Second)
	req.Retries(retry)
	req.Body([]byte(UrlEncode(data)))

	start := time.Now()

	ret, err := req.String()
	if err != nil {
		return "", errors.New(ToStr(time.Since(start).Seconds()) + " msg:" + err.Error())
	}
	return ret, nil
}

//PostJSON  Post  json
func PostJSON(ctx context.Context, apiurl string, data map[string]interface{}, second int, retry int) (string, error) {
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

	req.SetTimeout(1*time.Second, time.Duration(second)*time.Second)
	req.Retries(retry)
	req.Header("User-Agent", USER_AGENT)
	req.Header("Content-Type", "application/json")
	req.Body([]byte(JSONToStr(data)))

	start := time.Now()
	ret, err := req.String()
	if err != nil {
		return "", errors.New(ToStr(time.Since(start).Seconds()) + " msg:" + err.Error())
	}
	return ret, nil
}
