package main

import (
	"errors"
	"fmt"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/loadbalancer"
	"github.com/go-kit/kit/loadbalancer/static"
	"github.com/go-kit/kit/log"
	kitratelimit "github.com/go-kit/kit/ratelimit"
	httptransport "github.com/go-kit/kit/transport/http"
	jujuratelimit "github.com/juju/ratelimit"
	decoder "github.com/minchao/googlemaps-share-link-decoder"
	"github.com/sony/gobreaker"
	"golang.org/x/net/context"
	"io"
	"net/url"
	"strings"
	"time"
)

func proxyingMiddleware(proxyList string, ctx context.Context, logger log.Logger) ServiceMiddleware {
	if proxyList == "" {
		logger.Log("proxy_to", "none")
		return func(next decoder.Service) decoder.Service { return next }
	}

	proxies := split(proxyList)
	logger.Log("proxy_to", fmt.Sprint(proxies))

	return func(next decoder.Service) decoder.Service {
		var (
			qps         = 100 // max to each instance
			publisher   = static.NewPublisher(proxies, factory(ctx, qps), logger)
			lb          = loadbalancer.NewRoundRobin(publisher)
			maxAttempts = 3
			maxTime     = 10000 * time.Millisecond
			endpoint    = loadbalancer.Retry(maxAttempts, maxTime, lb)
		)
		return proxymw{ctx, endpoint, next}
	}
}

// proxymw implements decoder.Service, forwarding Decode requests to the provided endpoint
type proxymw struct {
	context.Context
	DecodeEndpoint endpoint.Endpoint
	decoder.Service
}

func (mw proxymw) Decode(req *decoder.Request) (rep *decoder.Response, err error) {
	resp, err := mw.DecodeEndpoint(mw.Context, req)
	if err != nil {
		return nil, err
	}

	r := resp.(response)
	if r.errorResponse != nil {
		return nil, errors.New(r.errorResponse.Error)
	}

	return r.Response, nil
}

func factory(ctx context.Context, qps int) loadbalancer.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		var e endpoint.Endpoint
		e = makeDecodeProxy(ctx, instance)
		e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)
		e = kitratelimit.NewTokenBucketLimiter(jujuratelimit.NewBucketWithRate(float64(qps), int64(qps)))(e)
		return e, nil, nil
	}
}

func makeDecodeProxy(ctx context.Context, instance string) endpoint.Endpoint {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		panic(err)
	}
	if u.Path == "" {
		u.Path = "/decode"
	}
	return httptransport.NewClient(
		"GET",
		u,
		encodeRequest,
		decodeResponse,
	).Endpoint()
}

func split(s string) []string {
	a := strings.Split(s, ",")
	for i := range a {
		a[i] = strings.TrimSpace(a[i])
	}
	return a
}
