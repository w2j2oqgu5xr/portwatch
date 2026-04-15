package history

import (
	"bytes"
	"io"
)

// newLineReader wraps a newline-delimited byte slice so that the standard
// json.Decoder can stream through each JSON object separated by newlines.
func newLineReader(data []byte) io.Reader {
	// Replace bare newlines between objects with a space so the decoder can
	// treat the whole file as a stream of top-level values.
	normalized := bytes.ReplaceAll(data, []byte("\n"), []byte(" "))
	return bytes.NewReader(normalized)
}
