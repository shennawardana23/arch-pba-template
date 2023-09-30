package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	var buf bytes.Buffer
	message := "this is message"
	Init()
	Logger.SetOutput(&buf)
	field := loggerField{}
	Info(context.Background(), message)
	json.Unmarshal(buf.Bytes(), &field)
	assert.Equal(t, message, field.Message)
	assert.EqualValues(t, severityInfo, field.Severity)
}
