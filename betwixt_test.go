package betwixt_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SimonRichardson/betwixt"
	"github.com/SimonRichardson/betwixt/pkg/output"
)

func TestPlaintext(t *testing.T) {
	t.Parallel()

	t.Run("basic", func(t *testing.T) {

		handler := http.DefaultServeMux
		handler.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("content-type", "application/json")
			bytes, _ := json.Marshal(map[string]string{
				"hello": "world",
			})
			w.WriteHeader(http.StatusOK)
			w.Write(bytes)
		})

		var (
			buffer  = new(bytes.Buffer)
			outputs = []betwixt.Output{
				output.NewPlaintext(output.MakeWriter(buffer)),
			}
			capture = betwixt.New(handler, outputs)
			server  = httptest.NewServer(capture)
		)

		request("GET", fmt.Sprintf("%s/hello", server.URL), nil, func(http.Header) {})
		request("GET", fmt.Sprintf("%s/hello?possible=1", server.URL), nil, func(http.Header) {})
		request("GET", fmt.Sprintf("%s/hello?possible=2", server.URL), nil, func(http.Header) {})
		request("GET", fmt.Sprintf("%s/hello?random=2", server.URL), nil, func(h http.Header) {
			h.Set("Header", "Value")
		})

		if err := capture.Output(); err != nil {
			t.Fatal(err)
		}

		body := `GET 200 - /hello
- Parameters:
 ・ possible 2 (optional)
 ・ random   2 (optional)
- Request Headers:
 ・ Accept-Encoding gzip
 ・ User-Agent      Go-http-client/1.1
 ・ Header          Value (optional)
- Response Headers:
 ・ Content-Type application/json
- Response Body:

  {"hello":"world"}
`
		if expected, actual := body, buffer.String(); expected != actual {
			t.Errorf("expected: \n%q\n, actual: \n%q\n", expected, actual)
		}
	})
}

func errored(s int) bool {
	return !(s == http.StatusOK || s == http.StatusCreated || s == http.StatusNoContent)
}

func request(reqType string, url string, payload []byte, fn func(http.Header)) []byte {
	client := &http.Client{}
	req, err := http.NewRequest(reqType, url, bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal(err)
	}

	req.ContentLength = int64(len(payload))
	fn(req.Header)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if errored(resp.StatusCode) {
		body, err := ioutil.ReadAll(resp.Body)
		fmt.Println(body, err)
		log.Fatal(fmt.Errorf("Request error: %d", resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()
	if err != nil {
		log.Fatal(err)
	}

	return body
}
