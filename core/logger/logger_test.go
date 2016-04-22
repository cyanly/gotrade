package logger_test

import (
	"fmt"
	"testing"

	l "github.com/cyanly/gotrade/core/logger"
)

func TestLogger_printf(t *testing.T) {
	l.Infof("logged in %s", "Chao")
}

func TestLogger_levels(t *testing.T) {
	l.Debug("uploading")
	l.Info("upload complete")
}

func TestLogger_WithFields(t *testing.T) {
	ctx := l.WithField("file", "1.png")
	ctx.Debug("uploading")
	ctx.Info("upload complete")

	ctx = l.WithField("payload", l.Fields{
		"file": "medium.png",
		"type": "image/png",
		"size": 1 << 20,
	})
	ctx.Debug("uploading")
	ctx.Info("upload complete")
}

func TestLogger_WithField(t *testing.T) {
	ctx := l.WithField("file", "3.png").WithField("user", "Chao")
	ctx.Debug("uploading")
	ctx.Info("upload complete")
}

func TestLogger_Trace_info(t *testing.T) {
	func() (err error) {
		defer l.WithField("file", "info.png").Trace("upload").Stop(&err)
		return nil
	}()
}

func TestLogger_Trace_error(t *testing.T) {
	func() (err error) {
		defer l.WithField("file", "error.png").Trace("upload").Stop(&err)
		return fmt.Errorf("boom")
	}()

}

func BenchmarkLogger_small(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l.Info("login")
	}
}

func BenchmarkLogger_medium(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l.WithField("payload", l.Fields{
			"file": "medium.png",
			"type": "image/png",
			"size": 1 << 20,
		}).Info("upload")
	}
}

//func BenchmarkLogger_large(b *testing.B) {
//	err := fmt.Errorf("boom")
//
//	for i := 0; i < b.N; i++ {
//		l.WithFields(l.Fields{
//			"file": "large.png",
//			"type": "image/png",
//			"size": 1 << 20,
//		}).WithFields(l.Fields{
//			"some":     "more",
//			"data":     "here",
//			"whatever": "blah blah",
//			"more":     "stuff",
//			"context":  "such useful",
//			"much":     "fun",
//		}).WithError(err).Error("upload failed")
//	}
//}
