package logger

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-stack/stack"
	"github.com/mochammadshenna/arch-pba-template/internal/util/json"
	"github.com/sirupsen/logrus"
)

type severity string

const (
	severityDebug    severity = "DEBUG"
	severityInfo     severity = "INFO"
	severityWarning  severity = "WARNING"
	severityError    severity = "ERROR"
	severityCritical severity = "CRITICAL"
	severityAlert    severity = "ALERT"
)

var levelsToSeverity = map[logrus.Level]severity{
	logrus.DebugLevel: severityDebug,
	logrus.InfoLevel:  severityInfo,
	logrus.WarnLevel:  severityWarning,
	logrus.ErrorLevel: severityError,
	logrus.FatalLevel: severityCritical,
	logrus.PanicLevel: severityAlert,
}

type sourceLocation struct {
	File     string `json:"file,omitempty"`
	Line     int    `json:"line,omitempty"`
	Function string `json:"function,omitempty"`
}

type entry map[string]interface{}

// Formatter implements Stackdriver formatting for logrus.
type Formatter struct {
	Service   string
	Version   string
	StackSkip []string
}

// Option lets you configure the Formatter.
type Option func(*Formatter)

// WithStackSkip lets you configure which packages should be skipped for locating the error.
func WithStackSkip(v string) Option {
	return func(f *Formatter) {
		f.StackSkip = append(f.StackSkip, v)
	}
}

// NewFormatter returns a new Formatter.
func NewFormatter(options ...Option) *Formatter {
	fmtr := Formatter{
		StackSkip: []string{
			"github.com/sirupsen/logrus",
		},
	}
	for _, option := range options {
		option(&fmtr)
	}
	return &fmtr
}

func (f *Formatter) errorOrigin() (stack.Call, error) {
	skip := func(pkg string) bool {
		for _, skip := range f.StackSkip {
			if pkg == skip {
				return true
			}
		}
		return false
	}

	// We start at 2 to skip this call and our caller's call.
	for i := 2; ; i++ {
		c := stack.Caller(i)
		// ErrNoFunc indicates we're over traversing the stack.
		if _, err := c.MarshalText(); err != nil {
			return stack.Call{}, nil
		}
		pkg := fmt.Sprintf("%+k", c)
		// Remove vendoring from package path.
		parts := strings.SplitN(pkg, "/vendor/", 2)
		pkg = parts[len(parts)-1]
		if !skip(pkg) {
			return c, nil
		}
	}
}

// Format formats a logrus entry according to the Stackdriver specifications.
func (f *Formatter) Format(e *logrus.Entry) ([]byte, error) {
	severity := levelsToSeverity[e.Level]

	ee := entry{
		LoggerField().Message:  e.Message,
		LoggerField().Severity: severity,
	}

	ee[LoggerField().Timestamp] = time.Now().UTC().Format(time.RFC3339)

	// Extract all custom field and push into log
	for k, v := range e.Data {
		ee[k] = v
	}

	switch severity {
	case severityError, severityCritical, severityAlert:

		// When using WithError(), the error is sent separately, but Error
		// Reporting expects it to be a part of the message so we append it
		// instead.
		if err, ok := e.Data["error"]; ok {
			ee["error"] = fmt.Sprintf("%s: %s", e.Message, err)
		} else {
			ee["error"] = e.Message
		}

		// Extract report location from call stack.
		if c, err := f.errorOrigin(); err == nil {
			lineNumber, _ := strconv.ParseInt(fmt.Sprintf("%d", c), 10, 64)

			ee[LoggerField().SourceLocation] = &sourceLocation{
				File:     fmt.Sprintf("%+s", c),
				Line:     int(lineNumber),
				Function: fmt.Sprintf("%n", c),
			}
		}
	}

	b, err := json.Marshal(ee)
	if err != nil {
		return nil, err
	}

	return append(b, '\n'), nil
}
