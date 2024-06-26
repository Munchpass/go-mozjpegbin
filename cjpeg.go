package mozjpegbin

import (
	"errors"
	"fmt"
	"image"
	"io"

	"github.com/Munchpass/go-mozjpegbin/embedbinwrapper"
)

// CJpeg wraps cjpeg tool from mozjpeg
type CJpeg struct {
	BinWrapper *embedbinwrapper.EmbedBinWrapper
	inputFile  string
	inputImage image.Image
	input      io.Reader
	outputFile string
	output     io.Writer
	quality    int
	optimize   bool
}

// NewCJpeg creates new CJpeg instance
func NewCJpeg() (*CJpeg, error) {
	binWrapper, err := createBinWrapper("cjpeg")
	if err != nil {
		return nil, fmt.Errorf("failed to create bin wrapper: %v", err)
	}

	bin := &CJpeg{
		BinWrapper: binWrapper,
		quality:    -1,
	}

	return bin, nil
}

// InputFile sets image file to convert.
// Input or InputImage called before will be ignored.
func (c *CJpeg) InputFile(file string) *CJpeg {
	c.input = nil
	c.inputImage = nil
	c.inputFile = file
	return c
}

// Input sets reader to convert.
// InputFile or InputImage called before will be ignored.
func (c *CJpeg) Input(reader io.Reader) *CJpeg {
	c.inputFile = ""
	c.inputImage = nil
	c.input = reader
	return c
}

// InputImage sets image to convert.
// InputFile or Input called before will be ignored.
func (c *CJpeg) InputImage(img image.Image) *CJpeg {
	c.inputFile = ""
	c.input = nil
	c.inputImage = img
	return c
}

// OutputFile specify the name of the output jpeg file.
// Output called before will be ignored.
func (c *CJpeg) OutputFile(file string) *CJpeg {
	c.output = nil
	c.outputFile = file
	return c
}

// Output specify writer to write jpeg file content.
// OutputFile called before will be ignored.
func (c *CJpeg) Output(writer io.Writer) *CJpeg {
	c.outputFile = ""
	c.output = writer
	return c
}

// Quality specify the compression factor for RGB channels between 0 and 100. The default is 75.
//
// A small factor produces a smaller file with lower quality. Best quality is achieved by using a value of 100.
func (c *CJpeg) Quality(quality uint) *CJpeg {
	if quality > 100 {
		quality = 100
	}

	c.quality = int(quality)
	return c
}

// Optimize perform optimization of entropy encoding parameters.
// Without this, default encoding parameters are used.
// Optimize usually makes the JPEG file a little smaller, but cjpeg runs somewhat slower and needs much more memory.
// Image quality and speed of decompression are unaffected by Optimize.
func (c *CJpeg) Optimize(optimize bool) *CJpeg {
	c.optimize = optimize
	return c
}

// Run starts cjpeg with specified parameters.
func (c *CJpeg) Run() error {
	defer c.BinWrapper.Reset()

	if c.quality > -1 {
		c.BinWrapper.Arg("-quality", fmt.Sprintf("%d", c.quality))
	}

	if c.optimize {
		c.BinWrapper.Arg("-optimize")
	}

	output, err := c.getOutput()

	if err != nil {
		return err
	}

	if output != "" {
		c.BinWrapper.Arg("-outfile", output)
	}

	err = c.setInput()

	if err != nil {
		return err
	}

	if c.output != nil {
		c.BinWrapper.SetStdOut(c.output)
	}

	err = c.BinWrapper.Run()

	if err != nil {
		return errors.New(err.Error() + ". " + string(c.BinWrapper.StdErr()))
	}

	return nil
}

// Version returns cjpeg version.
func (c *CJpeg) Version() (string, error) {
	return version(c.BinWrapper)
}

// Reset resets all parameters to default values
func (c *CJpeg) Reset() *CJpeg {
	c.quality = -1
	c.optimize = false
	return c
}

func (c *CJpeg) setInput() error {
	if c.input != nil {
		c.BinWrapper.StdIn(c.input)
	} else if c.inputImage != nil {
		r, err := createReaderFromImage(c.inputImage)

		if err != nil {
			return err
		}

		c.BinWrapper.StdIn(r)
	} else if c.inputFile != "" {
		c.BinWrapper.Arg(c.inputFile)
	} else {
		return errors.New("undefined input")
	}

	return nil
}

func (c *CJpeg) getOutput() (string, error) {
	if c.output != nil {
		return "", nil
	} else if c.outputFile != "" {
		return c.outputFile, nil
	} else {
		return "", errors.New("undefined output")
	}
}
