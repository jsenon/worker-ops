//Package restapi is a package to handle all rest api
package restapi

import (
	"context"
	"encoding/json"
	"net/http"
	"runtime"
	"strconv"

	opentracing "github.com/opentracing/opentracing-go"
	tracelog "github.com/opentracing/opentracing-go/log"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/jsenon/worker-ops/internal/generate"
	"github.com/jsenon/worker-ops/internal/slack"
)

// A healthCheckResponse respresent configuration of the application
// swagger:response healthCheckResponse
type healthCheckResponse struct {
	Status string `json:"status"`
}

// A wellknownResponse respresent configuration of the application
// swagger:response wellknownResponse
type wellknownResponse struct {
	Servicename        string `json:"Servicename"`
	Servicedescription string `json:"Servicedescription"`
	Version            string `json:"Version"`
	Versionfull        string `json:"Versionfull"`
	Revision           string `json:"Revision"`
	Branch             string `json:"Branch"`
	Builddate          string `json:"Builddate"`
	Swaggerdocurl      string `json:"Swaggerdocurl"`
	Healthzurl         string `json:"Healthzurl"`
	Metricurl          string `json:"Metricurl"`
	Endpoints          string `json:"Endpoints"`
}

//Fake metrics struct
//swagger:response someResponse
type someResponse struct {
}

// WellKnownFingerHandler will provide the information about the service.
func WellKnownFingerHandler(w http.ResponseWriter, _ *http.Request) {

	// swagger:route GET /.well-known wellknown
	//
	// Have Well known Info.
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http
	//
	//
	//     Responses:
	//       200: wellknownResponse
	span, _ := opentracing.StartSpanFromContext(context.Background(), "(*woker-ops).WellKnownFingerHandler")
	span.LogFields(
		tracelog.String("event", "Received REST /.well-known"))
	defer span.Finish()
	item := wellknownResponse{
		Servicename:        "",
		Servicedescription: "",
		Version:            "",
		Versionfull:        "",
		Revision:           "",
		Branch:             "",
		Builddate:          "",
		Swaggerdocurl:      "",
		Healthzurl:         "/healthz",
		Metricurl:          "/metrics",
		Endpoints:          "/"}
	data, err := json.Marshal(item)
	if err != nil {
		log.Error().Msgf("Error %s", err.Error())
		span.LogFields(tracelog.String("Error", err.Error()))
		runtime.Goexit()
	}
	writeJSONResponse(span, w, http.StatusOK, data)
}

//Health will provide the information about state of the service.
func Health(w http.ResponseWriter, _ *http.Request) {

	// swagger:route GET /healthz healthz
	//
	// Have Health Info.
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http
	//
	//
	//     Responses:
	//       200: healthCheckResponse
	span, _ := opentracing.StartSpanFromContext(context.Background(), "(*woker-ops).WellKnownFingerHandler")
	span.LogFields(
		tracelog.String("event", "Received REST /healthz"))
	defer span.Finish()
	data, err := json.Marshal(healthCheckResponse{Status: "UP"})
	if err != nil {
		log.Error().Msgf("Error %s", err.Error())
		span.LogFields(tracelog.String("Error", err.Error()))
		runtime.Goexit()
	}
	log.Debug().Msgf("Debug Marshall health: %v", data)
	writeJSONResponse(span, w, http.StatusOK, data)

}

//Report send json information of worker instance that run since 12h
func Report(w http.ResponseWriter, r *http.Request) {

	// swagger:route GET /report report
	//
	// Have the report of all worker node.
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http
	//
	//
	//     Responses:
	//       200: body:FullInstances
	log.Debug().Msgf("Sarting Reporting")
	ctx := r.Context()

	vartime := viper.GetInt("apibefore")

	myresult, err := generate.Launch(ctx, vartime)
	if err != nil {
		log.Error().Msgf("Error %s", err.Error())
	}
	log.Debug().Msgf("Worker before JSON transformation: %v", myresult)
	json, err := json.MarshalIndent(myresult, "", "    ")
	if err != nil {
		log.Error().Msgf("Error %s", err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(json)))

	w.Write(json)

}

//SendReport will send report to slack/mail
func SendReport(w http.ResponseWriter, r *http.Request) {

	// swagger:route GET /send sendreport
	//
	// Send the reports.
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http
	//
	//
	//     Responses:
	//       200: someResponse
	log.Debug().Msgf("Sarting Sending")
	ctx := r.Context()

	vartime := viper.GetInt("apibefore")

	myresult, err := generate.Launch(ctx, vartime)
	if err != nil {
		log.Error().Msgf("Error %s", err.Error())
	}

	log.Debug().Msgf("To slack transfo: %v", myresult)
	slack.Tomsg(ctx, myresult)

}

//Metric Fake for swagger
func metrics() {

	// swagger:route GET /metrics metrics
	//
	// Have prometheus metrics.
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http
	//
	//
	//     Responses:
	//       200: someResponse
}

// writeJsonResponse will convert response to json
func writeJSONResponse(ctx opentracing.Span, w http.ResponseWriter, status int, data []byte) {
	sp := opentracing.StartSpan(
		"(*worker-ops).writeJSONResponse",
		opentracing.ChildOf(ctx.Context()))
	defer sp.Finish()
	sp.LogFields(tracelog.String("event", "Write string to JSON"))
	defer sp.Finish()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(status)
	_, err := w.Write(data)
	if err != nil {
		log.Error().Msgf("Error %s", err.Error())
		sp.LogFields(tracelog.String("Error", err.Error()))
	}
}
