package output

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/SimonRichardson/betwixt/pkg/entry"
)

type Plaintext struct {
	w io.WriteCloser
}

func NewPlaintext(w io.WriteCloser) *Plaintext {
	return &Plaintext{w}
}

func (o Plaintext) Output(docs []entry.Document) error {
	for _, v := range docs {
		fmt.Fprintf(o.w, "%s %s - %s\n", v.Method.String(), v.Status.String(), v.URL.String())
		fmt.Fprintf(o.w, "- Parameters:\n")

		// Common
		writeMap(o.w, v.Params)

		fmt.Fprintln(o.w, "- Request Headers:")

		writeMap(o.w, v.ReqHeaders)

		if union := v.ReqBody.String(); len(union) > 0 {
			fmt.Fprintln(o.w, "- Request Body:")
			fmt.Fprintf(o.w, "\n  %s\n\n", union)
		}

		fmt.Fprintln(o.w, "- Response Headers:")

		writeMap(o.w, v.RespHeaders)

		if union := v.RespBody.String(); len(union) > 0 {
			fmt.Fprintln(o.w, "- Response Body:")
			fmt.Fprintf(o.w, "\n  %s\n", union)
		}
	}

	o.w.Close()

	return nil
}

func writeMap(w io.Writer, params *entry.Map) {
	writer := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)

	params.Union().Values.Walk(func(k string, v interface{}) {
		fmt.Fprintf(writer, "\t・\t%s\t%v\n", k, entry.ToStrings(v).Join())
	})
	for _, v := range params.Difference() {
		v.Values.Walk(func(k string, v interface{}) {
			fmt.Fprintf(writer, "\t・\t%s\t%v (optional)\n", k, entry.ToStrings(v).Join())
		})
	}
	writer.Flush()
}
