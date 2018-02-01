package betwixt

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/SimonRichardson/betwixt/pkg/entry"
)

// Betwixt is a struct that holds all the entries and outputs to be processed
type Betwixt struct {
	mutex   sync.Mutex
	entries []entry.Entry
	outputs []Output
	handler http.Handler
}

// New creates a Betwixt for possible outputs
func New(handler http.Handler, outputs []Output) *Betwixt {
	return &Betwixt{
		mutex:   sync.Mutex{},
		outputs: outputs,
		handler: handler,
	}
}

// ServeHTTP handles all the middleware for creating the documents
func (b *Betwixt) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writer := httptest.NewRecorder()
	b.handler.ServeHTTP(writer, r)
	writer.Flush()

	nativeHeaders := w.Header()
	for k, v := range writer.Header() {
		for _, v := range v {
			nativeHeaders.Add(k, v)
		}
	}
	w.WriteHeader(writer.Code)
	w.Write(writer.Body.Bytes())

	// Only handle successful status codes
	if status := writer.Code; status >= 200 && status < 300 {
		b.mutex.Lock()
		defer b.mutex.Unlock()

		b.entries = append(b.entries, entry.Entry{
			URL:         r.URL,
			Method:      r.Method,
			Status:      writer.Code,
			ReqHeaders:  r.Header,
			ReqBody:     readBody(r.Body),
			RespHeaders: writer.Header(),
			RespBody: func() []byte {
				return writer.Body.Bytes()
			},
		})
	}
}

// Output the results
func (b *Betwixt) Output() error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	grouped, err := group(b.entries)
	if err != nil {
		return err
	}

	for _, v := range b.outputs {
		if err := v.Output(grouped); err != nil {
			return err
		}
	}

	return nil
}

// Output defines an interface for consuming a document.
type Output interface {
	Output([]entry.Document) error
}

func readBody(read io.ReadCloser) func() []byte {
	defer read.Close()

	var buffer bytes.Buffer
	if _, err := buffer.ReadFrom(read); err != nil {
		return func() []byte {
			return make([]byte, 0, 0)
		}
	}

	return buffer.Bytes
}

func group(entries []entry.Entry) ([]entry.Document, error) {
	// Group according to the url and status code.
	groups := entry.Entries(entries).GroupBy(func(entry entry.Entry) string {
		url := fmt.Sprintf("%s/%s", entry.URL.Host, entry.NormalisePath())
		return fmt.Sprintf("%s-%s-%d", entry.Method, url, entry.Status)
	})

	// Loop through all the groups and find differences.
	return groups.Walk(func(entries entry.Entries) (entry.Document, error) {
		url := entries.URL()
		return entry.Document{
			URL:         url,
			Method:      entries.Method(),
			Status:      entries.Status(),
			Params:      entries.Params(),
			ReqHeaders:  entries.ReqHeaders(),
			ReqBody:     entries.ReqBody(),
			RespHeaders: entries.RespHeaders(),
			RespBody:    entries.RespBody(),
		}, nil
	})
}
