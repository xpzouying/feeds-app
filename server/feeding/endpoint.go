package feeding

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func MakeListFeedsEndpoint(s Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listFeedsRequest)
		feeds := s.ListFeeds(req.Page, req.Count)
		return listFeedsResponse{feeds, nil}, nil
	}
}
