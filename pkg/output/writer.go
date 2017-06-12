package output

import "io"

type wrap struct {
	writer io.Writer
	closer func() error
}

// MakeWriter creates a io.WriteCloser out of any io.Writer
func MakeWriter(w io.Writer) io.WriteCloser {
	return wrap{w, func() error { return nil }}
}

// MakeWriterCloser creates a io.WriterCloser
func MakeWriterCloser(w io.Writer, fn func() error) io.WriteCloser {
	return wrap{w, fn}
}

func (w wrap) Write(p []byte) (int, error) {
	return w.writer.Write(p)
}

func (w wrap) Close() error {
	return w.closer()
}
