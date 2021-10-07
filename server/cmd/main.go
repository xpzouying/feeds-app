package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/xpzouying/feeds-app/server/feeding"
	"github.com/xpzouying/feeds-app/server/repository"
)

func main() {
	var (
		httpAddr string
	)
	{
		flag.StringVar(&httpAddr, "http.addr", ":8080", "http address for server")
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

	var fs feeding.Service
	{
		fs = feeding.NewService(feedRepo)
		fs = feeding.LoggingMiddleware(log.With(logger, "component", "feeding"))(fs)
	}

	mux := http.NewServeMux()
	mux.Handle("/feeding/", feeding.MakeHandler(fs))

	mux.Handle("/metrics", promhttp.Handler())

	logger.Log("http.addr", httpAddr)
	logger.Log("finish", http.ListenAndServe(httpAddr, mux))
}
