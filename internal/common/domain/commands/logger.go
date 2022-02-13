package commands

import (
	"bytes"
	"context"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/rafaelcalleja/go-kit/internal/common/domain/middleware"
)

var (
	// LogEntryCtxKey is the context.Context key to store the request log entry.
	LogEntryCtxKey = &contextKey{"LogEntry"}
)

type DefaultLogger struct {
	formatter *DefaultLogFormatter
	NoColor   bool
}

func NewDefaultLogger() *DefaultLogger {
	color := true
	if runtime.GOOS == "windows" {
		color = false
	}

	return &DefaultLogger{
		formatter: &DefaultLogFormatter{
			Logger:  log.New(os.Stdout, "", log.LstdFlags),
			NoColor: !color,
		},
	}
}

func (d DefaultLogger) Handle(stack middleware.StackMiddleware, ctx context.Context, closure middleware.Closure) error {
	entry := d.formatter.NewLogEntry(ctx, stack, GetPipelineContext(ctx))

	t1 := time.Now()
	defer func() {
		entry.Write(200, time.Since(t1), nil)
	}()

	return stack.Next().Handle(stack, WithLogEntry(ctx, entry), closure)
}

func NewLoggerMiddleware() middleware.Middleware {
	return NewDefaultLogger()
}

// LogFormatter initiates the beginning of a new LogEntry per request.
// See DefaultLogFormatter for an example implementation.
type LogFormatter interface {
	NewLogEntry(ctx context.Context, stack middleware.StackMiddleware, pipelineContext PipelineContext) LogEntry
}

// LogEntry records the final log when a request completes.
// See defaultLogEntry for an example implementation.
type LogEntry interface {
	Write(status, elapsed time.Duration, extra interface{})
	Panic(v interface{}, stack []byte)
}

// GetLogEntry returns the in-context LogEntry for a request.
func GetLogEntry(ctx context.Context) LogEntry {
	entry, _ := ctx.Value(LogEntryCtxKey).(LogEntry)
	return entry
}

// WithLogEntry sets the in-context LogEntry for a request.
func WithLogEntry(ctx context.Context, entry LogEntry) context.Context {
	return context.WithValue(ctx, LogEntryCtxKey, entry)
}

// LoggerInterface accepts printing to stdlib logger or compatible logger.
type LoggerInterface interface {
	Print(v ...interface{})
}

// DefaultLogFormatter is a simple logger that implements a LogFormatter.
type DefaultLogFormatter struct {
	Logger  LoggerInterface
	NoColor bool
}

// NewLogEntry creates a new LogEntry for the command.
func (l *DefaultLogFormatter) NewLogEntry(ctx context.Context, stack middleware.StackMiddleware, pipelineContext PipelineContext) LogEntry {
	useColor := !l.NoColor
	entry := &defaultLogEntry{
		DefaultLogFormatter: l,
		pipelineContext:     pipelineContext,
		buf:                 &bytes.Buffer{},
		useColor:            useColor,
	}

	reqID := GetCommandId(ctx)
	if reqID != "" {
		cW(entry.buf, useColor, nYellow, "[correlation_id: %s] ", reqID)
	}
	cW(entry.buf, useColor, nCyan, "\"")
	cW(entry.buf, useColor, bMagenta, "%s ", pipelineContext.Command.Type())
	cW(entry.buf, useColor, bWhite, "%v ", pipelineContext.Command)

	return entry
}

type defaultLogEntry struct {
	*DefaultLogFormatter
	pipelineContext PipelineContext
	buf             *bytes.Buffer
	useColor        bool
}

func (l *defaultLogEntry) Write(status, elapsed time.Duration, extra interface{}) {
	switch {
	case status < 200:
		cW(l.buf, l.useColor, bBlue, "%03d", status)
	case status < 300:
		cW(l.buf, l.useColor, bGreen, "%03d", status)
	case status < 400:
		cW(l.buf, l.useColor, bCyan, "%03d", status)
	case status < 500:
		cW(l.buf, l.useColor, bYellow, "%03d", status)
	default:
		cW(l.buf, l.useColor, bRed, "%03d", status)
	}

	l.buf.WriteString(" in ")
	if elapsed < 500*time.Millisecond {
		cW(l.buf, l.useColor, nGreen, "%s", elapsed)
	} else if elapsed < 5*time.Second {
		cW(l.buf, l.useColor, nYellow, "%s", elapsed)
	} else {
		cW(l.buf, l.useColor, nRed, "%s", elapsed)
	}

	l.Logger.Print(l.buf.String())
}

func (l *defaultLogEntry) Panic(v interface{}, stack []byte) {
	PrintPrettyStack(v)
}
