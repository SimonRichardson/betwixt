package betwixt_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/SimonRichardson/betwixt"
	"github.com/SimonRichardson/betwixt/pkg/output"
)

func Example_Plaintext() {
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
		log.Fatal(err)
	}

	fmt.Println(buffer.String())
	// Output:
	// GET 200 - /hello
	// - Parameters:
	//  ・ possible 2 (optional)
	//  ・ random   2 (optional)
	// - Request Headers:
	//  ・ Accept-Encoding gzip
	//  ・ User-Agent      Go-http-client/1.1
	//  ・ Header          Value (optional)
	// - Response Headers:
	//  ・ Content-Type application/json
	// - Response Body:
	//
	//   {"hello":"world"}
}
