package main

import (
	"io"
	"os"

	"github.com/Pavel7004/WebShop/pkg/adapters/db/mongo"
	"github.com/Pavel7004/WebShop/pkg/adapters/http"
	"github.com/Pavel7004/WebShop/pkg/components/shop"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	closer := InitTracing()
	defer closer.Close()

	db := mongo.New()

	shop := shop.New(db)
	server := http.New(shop)

	log.Info().Msg("Starting server")
	if err := server.Run(); err != nil {
		log.Error().Err(err)
	}
}

func InitTracing() io.Closer {
	cfg := jaegercfg.Configuration{
		ServiceName: "WebShop",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}
	jMetricsFactory := metrics.NullFactory
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(nil),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		panic(err)
	}
	opentracing.SetGlobalTracer(tracer)
	return closer
}
