package report

import (
	"bufio"
	"io"
)

// Renderer writes a complete report to an arbitrary sink. It buffers the
// output so the many small fmt calls in the header and table hit the
// underlying writer only once.
type Renderer struct {
	bw *bufio.Writer
}

// NewRenderer wraps w with a buffered writer ready for report output.
func NewRenderer(w io.Writer) *Renderer {
	return &Renderer{bw: bufio.NewWriter(w)}
}

// WriteReport renders the full report and flushes the buffer. Any error
// writing to the underlying sink is returned.
func (r *Renderer) WriteReport(res *Result, m Meta) error {
	if err := Render(r.bw, res, m); err != nil {
		return err
	}
	return r.bw.Flush()
}
