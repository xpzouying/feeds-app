package feeding

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/xpzouying/feeds-app/server/feed"
)

func MakeHandler(set EndpointSet) http.Handler {

	r := mux.NewRouter()

	r.Methods(http.MethodGet).Path("/feeding/feeds").Handler(kithttp.NewServer(
		set.ListFeeds,
		decodeListFeedsRequest,
		encodeResponse,
	))

	r.Methods(http.MethodPost).Path("/feeding/feeds").Handler(kithttp.NewServer(
		set.PostFeed,
		decodePostFeedRequest,
		encodeResponse,
	))

	return r
}

func decodeListFeedsRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var (
		page, _  = strconv.Atoi(r.FormValue("page"))
		count, _ = strconv.Atoi(r.FormValue("count"))
	)

	if count == 0 {
		count = 10
	}

	request = listFeedsRequest{
		Page:  page,
		Count: count,
	}
	return
}

func decodePostFeedRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var (
		uid, _ = strconv.Atoi(r.FormValue("uid"))
		text   = r.FormValue("text")
	)

	request = postFeedRequest{
		Uid:  uid,
		Text: text,
	}
	return
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		return json.NewEncoder(w).Encode(map[string]interface{}{
			"error": e.error().Error(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

type listFeedsRequest struct {
	Page  int `json:"page"`
	Count int `json:"count"`
}

type listFeedsResponse struct {
	Feeds []feed.Feed `json:"feeds"`
	Err   error       `json:"error,omitempty"`
}

type postFeedRequest struct {
	Uid  int    `json:"uid"`
	Text string `json:"text"`
}

type postFeedResponse struct {
	Feed feed.Feed `json:"feed"`
	Err  error     `json:"err,omitempty"`
}
