package main

import (
	"fmt"
	"github.com/go-kit/kit/metrics"
	decoder "github.com/minchao/googlemaps-share-link-decoder"
	"time"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.TimeHistogram
	countResult    metrics.Histogram
	decoder.Service
}

func (mw instrumentingMiddleware) Decode(req *decoder.Request) (rep *decoder.Response, err error) {
	defer func(begin time.Time) {
		methodField := metrics.Field{Key: "method", Value: "decode"}
		errorField := metrics.Field{Key: "error", Value: fmt.Sprintf("%v", err)}
		mw.requestCount.With(methodField).With(errorField).Add(1)
		mw.requestLatency.With(methodField).With(errorField).Observe(time.Since(begin))
	}(time.Now())

	rep, err = mw.Service.Decode(req)
	return
}
