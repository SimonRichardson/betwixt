package entry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Score is a type alias for a score, which is modelled as float64
type Score float64

// Entry defines a raw http request.
type Entry struct {
	URL         *url.URL
	Method      string
	Status      int
	ReqHeaders  http.Header
	ReqBody     func() []byte
	RespHeaders http.Header
	RespBody    func() []byte
}

// NormalisePath attempts to normalise both a Host and Path in a sane way
func (e Entry) NormalisePath() string {
	path := fmt.Sprintf("%s%s", e.URL.Host, e.URL.Path)
	for k, v := range e.URL.Query() {
		if strings.Index(k, ":") == 0 {
			value := strings.Join(v, "")
			if index := strings.Index(path, value); index >= 0 {
				path = path[:index] + k + path[index+len(value):]
			}
		}
	}
	return path
}

// Entries is a type alias for a slice of Entry
type Entries []Entry

// GroupBy defines a way to group entires by a particular key
func (e Entries) GroupBy(fn func(Entry) string) GroupedEntries {
	groups := make(GroupedEntries, 0)
	for _, v := range e {
		p := fn(v)
		groups[p] = append(groups[p], v)
	}
	return groups
}

// URL returns a URL of all possible entries.
func (e Entries) URL() *URL {
	u := NewURL()
	for _, v := range e {
		u.Add(HostPath{
			Host: v.URL.Host,
			Path: v.NormalisePath(),
		})
	}
	return u
}

// Method returns a String of all possible methods.
func (e Entries) Method() *String {
	m := NewString()
	for _, v := range e {
		m.Add(v.Method)
	}
	return m
}

// Status returns a Status of all possible http status
func (e Entries) Status() *Status {
	s := NewStatus()
	for _, v := range e {
		s.Add(v.Status)
	}
	return s
}

// Params returns a Map of all possible http parameters
func (e Entries) Params() *Map {
	p := NewMap()
	for _, v := range e {
		var (
			values = make(ValuesPromoted, 0)
			url    = v.URL
		)
		for k, v := range url.Query() {
			bytes, _ := json.Marshal(v)
			values[k] = ValuePromoted{
				Value:    string(bytes),
				Promoted: isURLKey(url.Path, k, v),
			}
		}
		p.Add(values)
	}
	return p
}

// ReqHeaders returns a Map of all possible http request headers
func (e Entries) ReqHeaders() *Map {
	p := NewMap()
	for _, v := range e {
		values := make(ValuesPromoted, 0)
		for k, v := range v.ReqHeaders {
			bytes, _ := json.Marshal(v)
			values[k] = ValuePromoted{
				Value: string(bytes),
			}
		}
		p.Add(values)
	}
	return p
}

// ReqBody returns a String of all possible http request bodies
func (e Entries) ReqBody() *String {
	m := NewString()
	for _, v := range e {
		m.Add(string(v.ReqBody()))
	}
	return m
}

// RespHeaders returns a Map of all possible http response headers
func (e Entries) RespHeaders() *Map {
	p := NewMap()
	for _, v := range e {
		values := make(ValuesPromoted, 0)
		for k, v := range v.RespHeaders {
			bytes, _ := json.Marshal(v)
			values[k] = ValuePromoted{
				Value: string(bytes),
			}
		}
		p.Add(values)
	}
	return p
}

// RespBody returns a String of all possible http response bodies
func (e Entries) RespBody() *String {
	m := NewString()
	for _, v := range e {
		m.Add(string(v.RespBody()))
	}
	return m
}

// GroupedEntries allows the grouping of all entries for a specific key
type GroupedEntries map[string][]Entry

// Walk allows the walking over of each Entry returning a Document
func (g GroupedEntries) Walk(fn func(Entries) (Document, error)) ([]Document, error) {
	res := make([]Document, 0, len(g))
	for _, v := range g {
		o, err := fn(Entries(v))
		if err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil
}

// Strings is a type alias for a slice of strings
type Strings []string

// ToStrings attempts to convert a interface{} to a string if possible.
func ToStrings(x interface{}) Strings {
	switch t := x.(type) {
	case []interface{}:
		var res Strings
		for _, v := range []interface{}(t) {
			if x, ok := v.(string); ok {
				res = append(res, x)
			}
		}
		return res
	}
	return make(Strings, 0)
}

// Join joins all the strings in a csv format.
func (s Strings) Join() string {
	return strings.Join(s, ", ")
}

func isURLKey(path, param string, value []string) bool {
	if strings.Index(param, ":") == 0 {
		v := strings.Join(value, "")
		if index := strings.Index(path, v); index >= 0 {
			return true
		}
	}
	return false
}
