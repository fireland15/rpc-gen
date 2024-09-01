package writer

import (
	"io"

	"github.com/fireland15/rpc-gen/compiler"
)

type Writer struct {
	Config *WriterConfig
}

func (w *Writer) WriteClient(service *compiler.Service, client string, out io.Writer) error {

}
