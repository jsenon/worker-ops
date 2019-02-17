// Package rest Worker-Ops API.
//
// the purpose of this application is to provide an application
// to monitor Worker Node
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http
//     Host: localhost
//     BasePath: /
//     Version: 0.0.1
//     License: Apache 2.0 https://opensource.org/licenses/Apache-2.0
//     Contact: Julien SENON<julien.senon@delair.aero>
//
//
//     Produces:
//     - application/json
//
// swagger:meta
//Package rest launch server
package rest

import (
	"context"
	"net/http"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"

	intprometheus "github.com/jsenon/worker-ops/internal/prometheus"
	"github.com/jsenon/worker-ops/internal/restapi"
	pkgprometheus "github.com/jsenon/worker-ops/pkg/prometheus"

	tracelog "github.com/opentracing/opentracing-go/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//ServeRest will launch Gin Server and server routes
func ServeRest(ctx context.Context) {
	go func() {
		for {
			parent := opentracing.StartSpan("Metrics")
			defer parent.Finish()
			parent.LogFields(
				tracelog.String("event", "Prometheus metrics launch"))
			intprometheus.LauchRecord(ctx, parent)
			time.Sleep(1 * time.Second)
		}
	}()

	http.HandleFunc("/healthz", restapi.Health)
	http.HandleFunc("./.well-known", restapi.WellKnownFingerHandler)
	http.HandleFunc("/report", restapi.Report)
	http.HandleFunc("/send", restapi.SendReport)
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Error().Msgf("Error Starting server: %v", err.Error())
	}

}

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(pkgprometheus.WorkerNumber)
}
