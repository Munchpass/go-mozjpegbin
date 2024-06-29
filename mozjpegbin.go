package mozjpegbin

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Munchpass/go-mozjpegbin/embedbinwrapper"
)

//go:embed bin/*
var binariesFs embed.FS

func createBinWrapper(binaryName string) (*embedbinwrapper.EmbedBinWrapper, error) {
	b := embedbinwrapper.NewExecutableBinWrapper()
	switch runtime.GOOS {
	case "windows":
		binPath := fmt.Sprintf("bin/windows/%s", binaryName)
		ext := strings.ToLower(filepath.Ext(binPath))
		if ext != ".exe" {
			binPath += ".exe"
		}

		binary, err := binariesFs.ReadFile(binPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read embed binary: %s", err)
		}
		return b.Src(embedbinwrapper.NewSrc().Bin(binary).Os("win32")), nil
	case "linux":
		binary, err := binariesFs.ReadFile(fmt.Sprintf("bin/linux/%s", binaryName))
		if err != nil {
			return nil, fmt.Errorf("failed to read embed binary: %s", err)
		}
		return b.Src(embedbinwrapper.NewSrc().Bin(binary).Os("linux")), nil

	case "darwin":
		binary, err := binariesFs.ReadFile(fmt.Sprintf("bin/darwin/%s", binaryName))
		if err != nil {
			return nil, fmt.Errorf("failed to read embed binary: %s", err)
		}
		return b.Src(embedbinwrapper.NewSrc().Bin(binary).Os("darwin")), nil
	default:
		return nil, fmt.Errorf("unsupported OS %s", runtime.GOOS)
	}
}

func createReaderFromImage(img image.Image) (io.Reader, error) {
	var buffer bytes.Buffer
	err := jpeg.Encode(&buffer, img, &jpeg.Options{Quality: 100})
	return &buffer, err
}

func version(b *embedbinwrapper.EmbedBinWrapper) (string, error) {
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
