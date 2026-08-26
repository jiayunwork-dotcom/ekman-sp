package report

import (
	"bufio"
	"io"
)

type Renderer struct {
	bw *bufio.Writer
}

func NewRenderer(w io.Writer) *Renderer {
	return &Renderer{bw: bufio.NewWriter(w)}
}

func (r *Renderer) WriteReport(res *Result, m Meta) error {
	if err := Render(r.bw, res, m); err != nil {
		return err
	}
	return r.bw.Flush()
}
