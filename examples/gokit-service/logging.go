package main

import (
	"encoding/json"
	"github.com/go-kit/kit/log"
	decoder "github.com/minchao/googlemaps-share-link-decoder"
	"time"
)

func loggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next decoder.Service) decoder.Service {
		return logmw{logger, next}
	}
}

type logmw struct {
	logger log.Logger
	decoder.Service
}

func (mw logmw) Decode(req *decoder.Request) (rep *decoder.Response, err error) {
	defer func(begin time.Time) {
		var out string
		if rep != nil {
			o, _ := json.Marshal(rep)
			out = string(o)
		}

		_ = mw.logger.Log(
			"method", "decode",
			"input", req.URL,
			"out", out,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	rep, err = mw.Service.Decode(req)
	return
}
