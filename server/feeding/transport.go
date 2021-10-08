package feeding

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/xpzouying/feeds-app/server/feed"
)

func MakeHandler(endpoint endpoint.Endpoint) http.Handler {

	listFeedsHandler := kithttp.NewServer(
		endpoint,
		decodeListFeedsRequest,
		encodeResponse,
	)

	r := mux.NewRouter()
	r.Methods(http.MethodGet).Path("/feeding/feeds").Handler(listFeedsHandler)

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
