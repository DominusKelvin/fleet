package logging

import (
	"bytes"
	"context"
	"strings"
	"testing"

	kitlog "github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestLoggingErrs(t *testing.T) {
	setupTest := func() (*bytes.Buffer, kitlog.Logger, *LoggingContext, context.Context) {
		buf := new(bytes.Buffer)
		logger := kitlog.NewLogfmtLogger(buf)
		lc := &LoggingContext{}
		ctx := NewContext(context.Background(), lc)
		return buf, logger, lc, ctx
	}
	checkLogEnds := func(t *testing.T, logLine string, expected string) bool {
		return assert.True(t, strings.HasSuffix(strings.TrimSpace(logLine), expected), logLine)
	}

	t.Run("one error", func(t *testing.T) {
		buf, logger, lc, ctx := setupTest()

		WithErr(ctx, errors.Wrap(errors.New("AAAA"), "BLAH"))
		lc.Log(ctx, logger)
		logLine := buf.String()
		checkLogEnds(t, logLine, `err="BLAH: AAAA"`)
	})
	t.Run("two errors", func(t *testing.T) {
		buf, logger, lc, ctx := setupTest()

		WithErr(ctx, errors.Wrap(errors.New("AAAA"), "BLAH"))
		WithErr(ctx, errors.Wrap(errors.New("BBBB"), "FOO"))
		lc.Log(ctx, logger)
		logLine := buf.String()
		checkLogEnds(t, logLine, `err="BLAH: AAAA || FOO: BBBB"`)
	})
}
