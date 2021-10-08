package main

import (
	"flag"
	"net/http"
	"os"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/log"
	stdopentracing "github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	stdzipkin "github.com/openzipkin/zipkin-go"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/xpzouying/feeds-app/server/feed"
	"github.com/xpzouying/feeds-app/server/feeding"
	"github.com/xpzouying/feeds-app/server/repository"
)

func main() {
	var (
		httpAddr   string
		zipkinAddr string
	)
	{
		flag.StringVar(&httpAddr, "http.addr", ":8080", "http address for server")
		flag.StringVar(&zipkinAddr, "zipkin.addr", "", "address of zipkin, like http://localhost:9411/api/v2/spans")
	}
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
		logger = log.With(logger, "ts", log.DefaultTimestamp)
	}

	var (
		feedRepo = repository.NewFeedRepository()
	)

	var (
		tracer = newOpenTracer(zipkinAddr)
	)

	var (
		fs feeding.Service = newFeedingService(logger, feedRepo)
	)

	var (
		endpoint = feeding.MakeListFeedsEndpoint(fs)
	)

	{
		endpoint = opentracing.TraceServer(tracer, "FeedingService")(endpoint)
	}

	mux := http.NewServeMux()
	mux.Handle("/feeding/", feeding.MakeHandler(endpoint))

	mux.Handle("/metrics", promhttp.Handler())

	logger.Log("http.addr", httpAddr)
	logger.Log("finish", http.ListenAndServe(httpAddr, mux))
}

func newOpenTracer(zipkinAddr string) (openTracer stdopentracing.Tracer) {
	openTracer = stdopentracing.GlobalTracer()
	if len(zipkinAddr) == 0 {
		return
	}

	zipkinReporter := httpreporter.NewReporter(zipkinAddr)
	zipkinEndpoint, err := stdzipkin.NewEndpoint("feed-server", ":0")
	if err != nil {
		return
	}

	zipkinTracer, err := stdzipkin.NewTracer(zipkinReporter, stdzipkin.WithLocalEndpoint(zipkinEndpoint))
	if err != nil {
		return
	}

	openTracer = zipkinot.Wrap(zipkinTracer)
	return
}

func newFeedingService(logger log.Logger, feedRepo feed.Repository) (fs feeding.Service) {
	labelNames := []string{"method"}

	fs = feeding.NewService(feedRepo)

	fs = feeding.LoggingMiddleware(log.With(logger, "component", "feeding"))(fs)
	fs = feeding.InstrumentMiddleware(
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "feeding_service",
			Name:      "request_count",
			Help:      "Count of request",
		}, labelNames),
		kitprometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
			Namespace: "api",
			Subsystem: "feeding_service",
			Name:      "request_latency",
			Help:      "Latency of request",
		}, labelNames),
	)(fs)

	return
}
