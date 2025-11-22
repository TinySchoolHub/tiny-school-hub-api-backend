package log

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/rs/zerolog"
)

func TestDebugZerolog(t *testing.T) {
	buf := &bytes.Buffer{}
	zlog := zerolog.New(buf).With().Timestamp().Logger()
	zlog = zlog.Level(zerolog.InfoLevel)

	zlog.Info().Msg("direct test")

	output := buf.String()
	fmt.Printf("Buffer output: %q\n", output)
	fmt.Printf("Buffer length: %d\n", buf.Len())

	if buf.Len() == 0 {
		t.Error("Buffer is empty!")
	}
}
