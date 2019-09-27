package log_test

import (
	"fmt"
	"testing"

	. "github.com/goookie/library/log"
)

func TestLogger(t *testing.T) {
	err := fmt.Errorf("no password")
	Logger.Error("test",
		ErrorField(err),
		Field("key1", "value1"),
	)
}
