# betwixt

Betwixt is a documentation tool to help API integration with external sources,
whilst development is ongoing. By placing the betwixt middleware between serving
API and clients it's possible to use integration API tests to output up to date
documentation for all.

This middleware is not router aligned so it integrates best with the standard
library, otherwise adapters can be created for usage with external routers.

## Output

Betwixt comes with three output formats; plaintext, markdown and Apiary
flavoured markdown. Alternative outputs can be easily added if required (json,
xml etc).

### Plaintext

To create plaintext output middleware:

```go
handler := http.DefaultServeMux
handler.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("content-type", "application/json")
    bytes, _ := json.Marshal(map[string]string{
        "hello": "world",
    })
    w.WriteHeader(http.StatusOK)
    w.Write(bytes)
})

var buffer bytes.Buffer
outputs := []betwixt.Output{
    output.NewPlaintext(output.MakeWriter(&buffer)),
}
capture := betwixt.New(handler, outputs)
server  := httptest.NewServer(capture)

...

if err := capture.Output(); err != nil {
    log.Fatal(err)
}

fmt.Println(buffer.String())
```

Example output for plaintext:

```
 Output:
 GET 200 - /hello
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
```
