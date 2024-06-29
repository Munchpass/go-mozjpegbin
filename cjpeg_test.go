package mozjpegbin_test

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/Munchpass/go-mozjpegbin"
	"github.com/stretchr/testify/assert"
)

func init() {
	downloadFile("https://upload.wikimedia.org/wikipedia/commons/e/e3/Avola-Syracuse-Sicilia-Italy_-_Creative_Commons_by_gnuckx_%283858115914%29.jpg", "source.jpg")
}

func downloadFile(url, target string) {
	_, err := os.Stat(target)

	if err != nil {
		resp, err := http.Get(url)

		if err != nil {
			fmt.Printf("Error while downloading test image: %v\n", err)
			panic(err)
		}

		defer resp.Body.Close()

		f, err := os.Create(target)

		if err != nil {
			panic(err)
		}

		defer f.Close()

		_, err = io.Copy(f, resp.Body)

		if err != nil {
			panic(err)
		}
	}
}

func TestEncodeImage(t *testing.T) {
	c, err := mozjpegbin.NewCJpeg()
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	f, err := os.Open("source.jpg")
	assert.Nil(t, err)
	img, err := jpeg.Decode(f)
	assert.Nil(t, err)
	c.InputImage(img)
	c.OutputFile("target.jpg")
	err = c.Run()
	assert.Nil(t, err)
	validateJpg(t)
}

func TestEncodeImage2(t *testing.T) {

	const width, height = 256, 256

	// Create a colored image of the given width and height.
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.NRGBA{
				R: uint8((x + y) & 255),
				G: uint8((x + y) << 1 & 255),
				B: uint8((x + y) << 2 & 255),
				A: 255,
			})
		}
	}

	c, err := mozjpegbin.NewCJpeg()
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	c.InputImage(img)
	c.OutputFile("target.jpg")
	err = c.Run()
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	validateJpgImage(t, img)
}

func TestEncodeReader(t *testing.T) {
	c, err := mozjpegbin.NewCJpeg()
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	f, err := os.Open("source.jpg")
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	c.Input(f)
	c.OutputFile("target.jpg")
	err = c.Run()
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	validateJpg(t)
}

func TestEncodeFile(t *testing.T) {
	c, err := mozjpegbin.NewCJpeg()
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	c.Quality(100)
	c.Optimize(true)
	c.InputFile("source.jpg")
	c.OutputFile("target.jpg")
	err = c.Run()
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	validateJpg(t)
}

func TestEncodeWriter(t *testing.T) {
	f, err := os.Create("target.jpg")
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	defer f.Close()

	c, err := mozjpegbin.NewCJpeg()
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	c.InputFile("source.jpg")
	c.Output(f)
	err = c.Run()
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	f.Close()
	validateJpg(t)
}

func TestCJpegVersion(t *testing.T) {
	c, err := mozjpegbin.NewCJpeg()
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	v, err := c.Version()
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	assert.NotEmpty(t, v)

	t.Logf("version: %s\n", v)
}

func validateJpg(t *testing.T) {
	//defer os.Remove("target.jpg")
	fSource, err := os.Open("source.jpg")
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	imgSource, err := jpeg.Decode(fSource)
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	validateJpgImage(t, imgSource)
}

func validateJpgImage(t *testing.T, imgSource image.Image) {
	//defer os.Remove("target.jpg")
	fTarget, err := os.Open("target.jpg")
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	defer fTarget.Close()
	imgTarget, err := jpeg.Decode(fTarget)
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	assert.Equal(t, imgSource.Bounds(), imgTarget.Bounds())
}
