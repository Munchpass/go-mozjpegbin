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
	if runtime.GOARCH == "arm" ||
		(runtime.GOOS != "windows" && runtime.GOOS != "linux" && runtime.GOOS != "darwin") {
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

func createBinWrapper(binaryName string) (*binwrapper.BinWrapper, error) {
	b := binwrapper.NewBinWrapper().AutoExe()

	if !skipDownload {
		switch runtime.GOOS {
		case "windows":
			// TODO: convert this to use the pre-built windows version
			b.Src(
				binwrapper.NewSrc().
					URL("https://mozjpeg.codelove.de/bin/mozjpeg_3.1_x86.zip").
					Os("win32"))
			return b.Strip(2).Dest(dest), nil
		case "linux":
			b.Src(
				binwrapper.NewSrc().ExecPath(fmt.Sprintf("./bin/linux/%s", binaryName)).
					Os("linux"))
			return b, nil
		case "darwin":
			b.Src(
				binwrapper.NewSrc().ExecPath(fmt.Sprintf("./bin/macos/%s", binaryName)).
					Os("darwin"))
			return b, nil
		default:
			return nil, fmt.Errorf("unsupported OS %s", runtime.GOOS)
		}
	}

	return b.Strip(2).Dest(dest), nil
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
