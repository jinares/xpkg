package xlog

import (
	"context"
	"errors"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/transport"
	"io"
	"os"
	"strings"
	"time"
)

var (
	//envServiceName    = "JAEGER_SERVICE_NAME"
	stdCloser         io.Closer
	jaegerServiceName string
)

//关闭tracer
func Close() {
	if stdCloser != nil {
		stdCloser.Close()
	}
}

func GetServiceName() string {
	if jaegerServiceName == "" {
		jaegerServiceName = os.Getenv(envServiceName)
	}
	return jaegerServiceName
}

//从环境变量读取配置初始化
func InitFromEnv() (err error) {
	cfg, err := config.FromEnv()
	if err != nil {
		logrus.Errorf("Could not parse Jaeger env vars: %s", err.Error())
		return err
	}
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		logrus.Errorf("Could not initialize jaeger tracer: %s", err.Error())
		if closer != nil {
			closer.Close()
		}
		return err
	}
	opentracing.SetGlobalTracer(tracer)
	stdCloser = closer
	return nil
}

//只生成链式日志，无需jaeger服务，不采样，不上报
func InitOnlyTracingLog(serviceName string) error {
	if serviceName == "" {
		serviceName = os.Getenv(envServiceName)
		if serviceName == "" {
			return errors.New("ServiceName can not be empty")
		}
	}
	jaegerServiceName = serviceName
	tracer, closer := jaeger.NewTracer(
		serviceName,
		jaeger.NewConstSampler(false),
		jaeger.NewNullReporter())
	opentracing.SetGlobalTracer(tracer)
	stdCloser = closer
	return nil
}

// Config 配置文件
type Config struct {
	ServiceName         string  `json:"ServiceName" yaml:"ServiceName"`                 //服务名, 必须
	JaegerAgent         string  `json:"JaegerAgent" yaml:"JaegerAgent"`                 //jaeger-agent地址，可选： 默认：jaeger-agent:6831
	SamplerRateLimit    float64 `json:"SamplerRateLimit" yaml:"SamplerRateLimit"`       //<=0 永不上报，< 1 按固定概率上报，> 1 限制每秒最大上报数量，=1 总是上报
	QueueSize           int     `json:"QueueSize" yaml:"QueueSize"`                     //异步span队列最大大小，可选：默认100
	BufferFlushInterval int     `json:"BufferFlushInterval" yaml:"BufferFlushInterval"` //buff刷新间隔，可选：默认1  单位：秒
}

//从配置文件初始化
func InitFromConfig(c *Config) (err error) {
	if c.QueueSize <= 0 {
		c.QueueSize = 100
	}
	if c.BufferFlushInterval <= 0 {
		c.BufferFlushInterval = 1
	}
	cfg := &config.Configuration{
		ServiceName: c.ServiceName,
		Sampler:     &config.SamplerConfig{},
		Reporter: &config.ReporterConfig{
			LogSpans:            false,
			QueueSize:           c.QueueSize,
			BufferFlushInterval: time.Duration(c.BufferFlushInterval) * time.Second,
		},
	}
	//永不上报
	if c.SamplerRateLimit <= 0 {
		cfg.Sampler.Type = jaeger.SamplerTypeConst
		cfg.Sampler.Param = 0
	} else if c.SamplerRateLimit < 1 {
		//按固定概率上报
		cfg.Sampler.Type = jaeger.SamplerTypeProbabilistic
		cfg.Sampler.Param = c.SamplerRateLimit
	} else if c.SamplerRateLimit > 1 {
		//每秒最大上报数量
		cfg.Sampler.Type = jaeger.SamplerTypeRateLimiting
		cfg.Sampler.Param = c.SamplerRateLimit
	} else {
		//总是上报
		cfg.Sampler.Type = jaeger.SamplerTypeConst
		cfg.Sampler.Param = 1
	}
	var opts []config.Option
	if c.JaegerAgent != "" {
		var sender jaeger.Transport
		if strings.HasPrefix(c.JaegerAgent, "http://") {
			sender = transport.NewHTTPTransport(
				c.JaegerAgent,
				transport.HTTPBatchSize(1),
			)
		} else {
			if s, err := jaeger.NewUDPTransport(c.JaegerAgent, 0); err != nil {
			} else {
				sender = s
			}
		}
		reporter := config.Reporter(jaeger.NewRemoteReporter(
			sender,
		))
		opts = append(opts, reporter)
	}
	tracer, closer, err := cfg.NewTracer(opts...)
	if err == nil {
		opentracing.SetGlobalTracer(tracer)
		stdCloser = closer
		jaegerServiceName = cfg.ServiceName
	}
	return err
}

func StartContext(ctx context.Context, operationName string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	tracer := opentracing.GlobalTracer()
	_, ctx = opentracing.StartSpanFromContextWithTracer(ctx, tracer, operationName)

	return ctx
}
func Carriers(ctx context.Context, op ...string) opentracing.TextMapCarrier {
	//df := ""
	//if len(op) > 0 {
	//	df = strings.Join(op, "")
	//}
	span := opentracing.SpanFromContext(ctx)
	carriers := opentracing.TextMapCarrier{}

	if span == nil {
		return carriers
	}
	err := span.Tracer().Inject(span.Context(), opentracing.TextMap, carriers)

	if err != nil {
		fmt.Println(fmt.Sprintf("carriers-inject-fail:%v", errString(err)))
		return opentracing.TextMapCarrier{}
	}
	return carriers
}

func ContextByCarrier(carrier opentracing.TextMapCarrier, operationName string, tag ...string) (ctx context.Context) {
	//defer func() {
	//	if TraceID(ctx) == "" {
	//		ctx = StartContext(context.Background(), operationName)
	//	}
	//}()
	tracer := opentracing.GlobalTracer()
	extractedContext, err := tracer.Extract(
		opentracing.TextMap,
		carrier,
	)
	if err != nil && err != opentracing.ErrSpanContextNotFound {
		fmt.Println("build-context-by-carrier-fail:", errString(err))
		return context.Background()
	} else {

		span := tracer.StartSpan(
			operationName,
			opentracing.ChildOf(extractedContext),
		)
		if span == nil {
			return context.Background()
		}
		if len(tag) > 0 {
			for _, value := range tag {
				ext.MessageBusDestination.Set(span, value)
			}
		}
		defer span.Finish()
		return opentracing.ContextWithSpan(context.Background(), span)
	}
}
