package mozjpegbin_test

import (
	"os"
	"testing"

	"github.com/Munchpass/go-mozjpegbin"
	"github.com/stretchr/testify/assert"
)

func TestJpegTranReader(t *testing.T) {
	c, err := mozjpegbin.NewJpegTran()
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

func TestJpegTranFile(t *testing.T) {
	c, err := mozjpegbin.NewJpegTran()
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	c.InputFile("source.jpg")
	c.OutputFile("target.jpg")
	err = c.Run()
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	validateJpg(t)
}

func TestJpegTranCrop(t *testing.T) {
	c, err := mozjpegbin.NewJpegTran()
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	c.Crop(500, 500, 100, 100)
	c.InputFile("source.jpg")
	c.OutputFile("target.jpg")
	err = c.Run()
	if !assert.Nil(t, err) {
		t.FailNow()
	}
}

func TestJpegTranWriter(t *testing.T) {
	f, err := os.Create("target.jpg")
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	defer f.Close()

	c, err := mozjpegbin.NewJpegTran()
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

func TestJpegTranVersion(t *testing.T) {
	c, err := mozjpegbin.NewJpegTran()
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	v, err := c.Version()
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	assert.NotZero(t, v)
}
