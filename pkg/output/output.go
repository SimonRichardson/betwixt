package output

import "github.com/SimonRichardson/betwixt/pkg/entry"

// Output defines an interface for consuming a document.
type Output interface {
	Output([]entry.Document) error
}
