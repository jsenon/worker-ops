//Package slack internal: Get url webhook, format message
package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/jsenon/worker-ops/internal/generate"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

//Geturl Retrieve Url form env
// TODO: Retrieve env webhook on consul
func geturl() (url string) {
	url = os.Getenv("SLACK_URL")
	log.Debug().Msgf("URL: %v", url)
	return url
}

//Tomsg Transform message
func Tomsg(ctx context.Context, sp opentracing.Span, msg []generate.FullInstances) {
	span := opentracing.StartSpan(
		"(*worker-ops).Tomsg",
		opentracing.ChildOf(sp.Context()))
	defer span.Finish()
	var slackmsg string

	url := geturl()

	for _, n := range msg {

		slackmsg = "Worker still running more than " + viper.GetString("apibefore") + " hour(s) in " + n.Env + " environment: \n"

		for _, o := range n.Instances {
			slackmsg = slackmsg + " Name: " + o.Name + " DNS: " + o.Dnsname + " on region: " + o.Region + " Started since UTC: " + o.Launchtime.String() + "\n"
		}

	}

	log.Debug().Msgf("Slack message reformated: %v", slackmsg)

	values := map[string]string{"text": slackmsg}

	b, err := json.Marshal(values)
	if err != nil {
		log.Error().Msgf("Error marshal: %v", err)
	}
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	httpclient := &http.Client{Transport: tr}
	rs, err := httpclient.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		panic(err)
	}
	defer rs.Body.Close() // nolint: errcheck

}
