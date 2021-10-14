package middleware

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/tracing/opentracing"
	stdopentracing "github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter/http"
)

func WithOpenTracer(ep endpoint.Endpoint, tracer stdopentracing.Tracer, traceName string) endpoint.Endpoint {

	return opentracing.TraceServer(tracer, traceName)(ep)
}

func NewOpenTracer(serviceName string, zipkinAddr string) (openTracer stdopentracing.Tracer) {
	openTracer = stdopentracing.GlobalTracer()
	if len(zipkinAddr) == 0 {
		return
	}

	if zipkinTracer, err := newZipkinOpenTracer(serviceName, zipkinAddr); err != nil {
		return
	} else {
		openTracer = zipkinTracer
	}
	return
}

func newZipkinOpenTracer(serviceName string, zipkinAddr string) (stdopentracing.Tracer, error) {
	zipkinReporter := http.NewReporter(zipkinAddr)
	zipkinEndpoint, err := stdzipkin.NewEndpoint(serviceName, ":0")
	if err != nil {
		return nil, err
	}

	zipkinTracer, err := stdzipkin.NewTracer(
		zipkinReporter,
		stdzipkin.WithLocalEndpoint(zipkinEndpoint),
	)
	if err != nil {
		return nil, err
	}

	return zipkinot.Wrap(zipkinTracer), nil
}
