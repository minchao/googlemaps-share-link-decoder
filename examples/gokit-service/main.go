package main

import (
	"flag"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	decoder "github.com/minchao/googlemaps-share-link-decoder"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"time"
)

func main() {
	var (
		listen = flag.String("listen", ":8080", "HTTP listen address")
		proxy  = flag.String("proxy", "", "Optional comma-separated list of URLs to proxy decode results")
	)
	flag.Parse()

	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.NewContext(logger).With("listen", *listen).With("caller", log.DefaultCaller)

	ctx := context.Background()

	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounter(stdprometheus.CounterOpts{
		Namespace: "example",
		Subsystem: "share_link_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := metrics.NewTimeHistogram(time.Microsecond, kitprometheus.NewSummary(stdprometheus.SummaryOpts{
		Namespace: "example",
		Subsystem: "share_link_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys))
	countResult := kitprometheus.NewSummary(stdprometheus.SummaryOpts{
		Namespace: "example",
		Subsystem: "share_link_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{})

	var svc decoder.Service
	svc = decoder.ShareLinkService{}
	svc = proxyingMiddleware(*proxy, ctx, logger)(svc)
	svc = loggingMiddleware(logger)(svc)
	svc = instrumentingMiddleware(requestCount, requestLatency, countResult)(svc)

	decodeHandler := httptransport.NewServer(
		ctx,
		makeDecodeEndpoint(svc),
		decodeRequest,
		encodeResponse,
	)

	http.Handle("/decode", decodeHandler)
	http.Handle("/metrics", stdprometheus.Handler())
	_ = logger.Log("msg", "HTTP", "addr", *listen)
	_ = logger.Log("err", http.ListenAndServe(*listen, nil))
}
