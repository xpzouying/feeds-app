package feeding

import (
	"context"
	"encoding/json"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/xpzouying/feeds-app/server/feed"
)

func MakeHandler(svc Service) http.Handler {

	listFeedsHandler := kithttp.NewServer(
		makeListFeedsEndpoint(svc),
		decodeListFeedsRequest,
		encodeResponse,
	)

	r := mux.NewRouter()
	r.Methods(http.MethodPost).Path("/feeding/feeds").Handler(listFeedsHandler)

	return r
}

func decodeListFeedsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request listFeedsRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	return request, err
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		return json.NewEncoder(w).Encode(map[string]interface{}{
			"error": e.error().Error(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
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
