package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/xpzouying/feeds-app/server/feeding"
	"github.com/xpzouying/feeds-app/server/middleware"
)

func main() {
	var (
		cfgPath string
	)
	{
		flag.StringVar(&cfgPath, "conf", "config.yaml", "config file for server")
	}
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
		logger = log.With(logger, "ts", log.DefaultTimestamp)
	}

	config, err := newConfigFromFile(cfgPath)
	if err != nil {
		logger.Log("error", err.Error())
		os.Exit(1)
	}

	fmt.Printf("config: %+v\n", config)

	repoSet, err := newRepositorySet(config.DB)
	if err != nil {
		logger.Log("error", err.Error())
		os.Exit(1)
	}

	var (
		fs feeding.Service = makeFeedingService(logger, repoSet.feedRepo, repoSet.userRepo)

		tracer = middleware.NewOpenTracer("feed-server", config.Tracer.Address)
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

	logger.Log("http.addr", config.Server.HTTPAddr)
	logger.Log("finish", http.ListenAndServe(config.Server.HTTPAddr, mux))
}
