package main

import (
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/log"
	stdprometheus "github.com/prometheus/client_golang/prometheus"

	"github.com/xpzouying/feeds-app/server/feed"
	"github.com/xpzouying/feeds-app/server/feeding"
	"github.com/xpzouying/feeds-app/server/repository/inmem"
	"github.com/xpzouying/feeds-app/server/repository/sqlrepo"
	"github.com/xpzouying/feeds-app/server/user"
)

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

type repositorySet struct {
	feedRepo feed.Repository
	userRepo user.Repository
}

func newRepositorySet(cfg dbConfig) (set repositorySet, err error) {

	if cfg.UseMem {
		set.feedRepo = inmem.NewFeedRepository()
		set.userRepo = inmem.NewUserRepository()
		return
	}

	option := sqlrepo.Option{
		User:     cfg.User,
		Password: cfg.Password,
		Host:     cfg.Host,
		Port:     cfg.Port,
	}

	if set.feedRepo, err = sqlrepo.NewFeedRepo(option); err != nil {
		return
	}
	set.userRepo = inmem.NewUserRepository()

	return
}
