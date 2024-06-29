package mozjpegbin

import (
	"fmt"
	"image"
	"io"
)

// Options to use with encoder
type Options struct {
	Quality  uint
	Optimize bool
}

// Encode encodes image.Image into jpeg using cjpeg.
func Encode(w io.Writer, m image.Image, o *Options) error {
	cjpeg, err := NewCJpeg()
	if err != nil {
		return fmt.Errorf("NewCJpeg failed: %v", err)
	}

	if o != nil {
		cjpeg.Quality(o.Quality)
		cjpeg.Optimize(o.Optimize)
	}

	return cjpeg.InputImage(m).Output(w).Run()
}
