//Package slack Get url, publish on slack
package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

//publish will send text to slack endpoint provided by webhook URL (URL contain the Token)
func publish(msg string, url string) {

	values := map[string]string{"text": msg}
	b, err := json.Marshal(values)
	if err != nil {
		fmt.Println(err)
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
