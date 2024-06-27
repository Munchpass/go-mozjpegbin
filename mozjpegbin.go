package mozjpegbin

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"runtime"
	"strings"

	"github.com/nickalie/go-binwrapper"
)

var skipDownload bool
var dest = "vendor/mozjpeg"

func init() {
	if runtime.GOARCH == "arm" || (runtime.GOOS != "windows" && runtime.GOOS != "linux") {
		SkipDownload()
	}
}

// SkipDownload skips binary download.
func SkipDownload() {
	skipDownload = true
	dest = ""
}

// Dest sets directory to download mozjpeg binaries or where to look for them if SkipDownload is used. Default is "vendor/mozjpeg"
func Dest(value string) {
	dest = value
}

func createBinWrapper(binaryName string) *binwrapper.BinWrapper {
	b := binwrapper.NewBinWrapper().AutoExe()

	if !skipDownload {
		if runtime.GOOS == "windows" {
			b.Src(
				binwrapper.NewSrc().
					URL("https://mozjpeg.codelove.de/bin/mozjpeg_3.1_x86.zip").
					Os("win32"))
		} else if runtime.GOOS == "linux" {
			b.Src(
				binwrapper.NewSrc().ExecPath(fmt.Sprintf("./bin/linux/%s", binaryName)).
					Os("linux"))
			return b
		}
	}

	return b.Strip(2).Dest(dest)
}

func createReaderFromImage(img image.Image) (io.Reader, error) {
	var buffer bytes.Buffer
	err := jpeg.Encode(&buffer, img, &jpeg.Options{Quality: 100})
	return &buffer, err
}

func version(b *binwrapper.BinWrapper) (string, error) {
	b.Reset()
	err := b.Run("-version")

	if err != nil {
		return "", err
	}

	v := string(b.StdErr())
	v = strings.Replace(v, "\n", "", -1)
	v = strings.Replace(v, "\r", "", -1)
	return v, nil
}
