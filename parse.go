package betwixt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/SimonRichardson/betwixt/pkg/output"
)

// Parse a string to a possible set of outputs
func Parse(value string) ([]Output, error) {
	var res []Output
	for _, v := range strings.Split(strings.ToLower(value), ";") {
		switch parts := strings.Split(v, ","); parts[0] {
		case "plaintext":
			out, err := getOutput(parts)
			if err != nil {
				return []Output{}, err
			}
			res = append(res, output.NewPlaintext(out))
		case "markdown":
			out, err := getOutput(parts)
			if err != nil {
				return []Output{}, err
			}
			res = append(res, output.NewMarkdown(out, getMarkdownOptions(parts)))
		}
	}
	return res, nil
}

func getOutput(parts []string) (io.WriteCloser, error) {
	if len(parts) < 2 {
		return os.Stdout, nil
	}

	switch value := strings.Split(parts[1], ":"); value[0] {
	case "stdout":
		return os.Stdout, nil
	case "file":
		if len(value) == 2 {
			path, err := filepath.Abs(value[1])
			if err != nil {
				return nil, err
			}
			file, err := os.Create(path)
			if err != nil {
				return nil, err
			}

			buf := bufio.NewWriter(file)
			return output.MakeWriterCloser(
				buf,
				func() error {
					buf.Flush()
					return file.Close()
				},
			), nil
		}
	}
	return nil, fmt.Errorf("no valid output found")
}

func getMarkdownOptions(parts []string) output.Options {
	if len(parts) < 3 {
		return output.Options{}
	}

	switch value := parts[2]; value {
	case "apiary":
		name := "Apiary"
		if len(parts) > 3 {
			name = parts[3]
		}
		return output.NewApiaryOptions(name)
	}

	return output.Options{}
}
