/*
Copyright © 2024 Kovalev Pavel

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package main

import (
	"io"
	"os"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"

	"github.com/Pavel7004/WebShop/pkg/adapters/db/mongo"
	"github.com/Pavel7004/WebShop/pkg/adapters/http"
	"github.com/Pavel7004/WebShop/pkg/components/shop"
	"github.com/Pavel7004/WebShop/pkg/infra/config"
)

func main() {
	closer := initTracing()
	defer closer.Close()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	cfg, err := config.Get()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get config")
		return
	}

	db := mongo.New(cfg)

	shop := shop.New(db)
	server := http.New(shop, cfg)

	log.Info().Msg("Starting server")
	if err := server.Run(); err != nil {
		log.Error().Err(err).Msg("Server error")
	}
}

func initTracing() io.Closer {
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
