package main

import (
	"net/http"
	"os"

	"github.com/go-kit/log"

	"github.com/xpzouying/feeds-app/server/feeding"
	"github.com/xpzouying/feeds-app/server/repository"
)

func main() {

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
	}

	mux := http.NewServeMux()
	mux.Handle("/feeding/", feeding.MakeHandler(fs))

	logger.Log("http.addr", ":8080")
	logger.Log("finish", http.ListenAndServe(":8080", mux))
}
