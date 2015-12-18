package main

import (
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	decoder "github.com/minchao/googlemaps-share-link-decoder"
	"golang.org/x/net/context"
	"net/http"
)

func makeDecodeEndpoint(svc decoder.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(decoder.Request)
		r, err := svc.Decode(&req)
		if err != nil {
			return errorResponse{err.Error()}, nil
		}
		return r, nil
	}
}

func decodeRequest(r *http.Request) (interface{}, error) {
	var request decoder.Request
	json.NewDecoder(r.Body).Decode(&request)
	return request, nil
}

func encodeResponse(w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

type errorResponse struct {
	URL string `json:"error"`
}
