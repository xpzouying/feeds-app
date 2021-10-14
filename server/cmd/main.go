package main

import (
	"flag"
	"net/http"
	"os"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/log"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/xpzouying/feeds-app/server/feed"
	"github.com/xpzouying/feeds-app/server/feeding"
	"github.com/xpzouying/feeds-app/server/middleware"
	"github.com/xpzouying/feeds-app/server/repository"
	"github.com/xpzouying/feeds-app/server/user"
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
		userRepo = repository.NewUserRepository()

		fs feeding.Service = makeFeedingService(logger, feedRepo, userRepo)

		tracer = middleware.NewOpenTracer("feed-server", zipkinAddr)
	)

	var feedingEndpoints feeding.EndpointSet
	{
		listFeedEndpoint := feeding.MakeListFeedsEndpoint(fs)
		listFeedEndpoint = middleware.WithRateLimiter(listFeedEndpoint)
		listFeedEndpoint = middleware.WithOpenTracer(listFeedEndpoint, tracer, "FeedingService.ListFeeds")
		feedingEndpoints.ListFeeds = listFeedEndpoint

		postFeedEndpoint := feeding.MakePostFeedEndpoint(fs)
		postFeedEndpoint = middleware.WithGoBreaker(postFeedEndpoint, "FeedingService.PostFeed.Breader")
		postFeedEndpoint = middleware.WithOpenTracer(postFeedEndpoint, tracer, "FeedingService.PostFeed")
		feedingEndpoints.PostFeed = postFeedEndpoint
	}

	mux := http.NewServeMux()
	mux.Handle("/feeding/", feeding.MakeHandler(feedingEndpoints))

	mux.Handle("/metrics", promhttp.Handler())

	logger.Log("http.addr", httpAddr)
	logger.Log("finish", http.ListenAndServe(httpAddr, mux))
}

func makeFeedingService(logger log.Logger, feedRepo feed.Repository, userRepo user.Repository) (fs feeding.Service) {
	labelNames := []string{"method"}

	fs = feeding.NewService(feedRepo, userRepo)

	fs = feeding.WithLoggingMiddleware(log.With(logger, "component", "feeding"))(fs)
	fs = feeding.WithInstrumentMiddleware(
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
