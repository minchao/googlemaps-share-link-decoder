package main

import (
	"bytes"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	decoder "github.com/minchao/googlemaps-share-link-decoder"
	"golang.org/x/net/context"
	"io/ioutil"
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

func encodeRequest(r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func decodeRequest(r *http.Request) (interface{}, error) {
	var request decoder.Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func decodeResponse(r *http.Response) (interface{}, error) {
	var response response
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

type errorResponse struct {
	Error string `json:"error"`
}

type response struct {
	*decoder.Response
	*errorResponse
}
