package feeding

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type EndpointSet struct {
	ListFeeds endpoint.Endpoint
	PostFeed  endpoint.Endpoint
}

func MakeListFeedsEndpoint(s Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listFeedsRequest)
		feeds := s.ListFeeds(req.Page, req.Count)
		return listFeedsResponse{feeds, nil}, nil
	}
}

func MakePostFeedEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postFeedRequest)
		feed, err := s.PostFeed(req.Uid, req.Text)
		return postFeedResponse{feed, err}, nil
	}
}
