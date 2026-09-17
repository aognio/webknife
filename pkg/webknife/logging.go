package webknife

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/webknife/webknife/internal/events"
)

type LogFormat string

const (
	LogFormatText LogFormat = "text"
	LogFormatJSON LogFormat = "json"
)

type LogAdapter struct {
	format LogFormat
	output io.Writer
}

func NewLogAdapter(format LogFormat, output io.Writer) *LogAdapter {
	return &LogAdapter{
		format: format,
		output: output,
	}
}

func (a *LogAdapter) Handle(event events.Event) {
	switch e := event.(type) {
	case events.RequestCompleted:
		a.logRequest(e)
	case events.ServerStarted:
		a.logServerStarted(e)
	case events.ServerStopped:
		a.logServerStopped(e)
	case events.AuthenticationFailed:
		a.logAuthFailed(e)
	}
}

func (a *LogAdapter) logServerStarted(e events.ServerStarted) {
	ts := time.Now().Format(time.RFC3339)
	mode := "HTTP"
	if e.TLS {
		mode = "HTTPS"
	}
	log.SetOutput(a.output)
	log.Printf("%s %s server listening on %s", ts, mode, e.Addr)
}

func (a *LogAdapter) logServerStopped(e events.ServerStopped) {
	ts := time.Now().Format(time.RFC3339)
	log.SetOutput(a.output)
	log.Printf("%s server stopped on %s", ts, e.Addr)
}

func (a *LogAdapter) logAuthFailed(e events.AuthenticationFailed) {
	ts := time.Now().Format(time.RFC3339)
	log.SetOutput(a.output)
	log.Printf("%s authentication failed user=%s path=%s", ts, e.User, e.Path)
}

func (a *LogAdapter) logRequest(e events.RequestCompleted) {
	sizeStr := formatBytes(e.Size)
	remote := e.RemoteAddr
	if a.format == LogFormatJSON {
		a.logJSON(e, sizeStr, remote)
	} else {
		a.logText(e, sizeStr, remote)
	}
}

func formatBytes(b int64) string {
	switch {
	case b >= 1<<20:
		return fmt.Sprintf("%.1fMB", float64(b)/(1<<20))
	case b >= 1<<10:
		return fmt.Sprintf("%.1fKB", float64(b)/(1<<10))
	default:
		return fmt.Sprintf("%dB", b)
	}
}

func (a *LogAdapter) logText(e events.RequestCompleted, sizeStr, remote string) {
	ts := e.Timestamp.Format(time.RFC3339)
	log.SetOutput(a.output)
	log.Printf("%s %s %s %d %s %s %s",
		ts, e.Method, e.Path, e.StatusCode, sizeStr, e.Duration.Round(time.Microsecond), remote)
}

type jsonEntry struct {
	Timestamp string `json:"timestamp"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Status    int    `json:"status"`
	Size      string `json:"size"`
	Duration  string `json:"duration"`
	Remote    string `json:"remote"`
}

func (a *LogAdapter) logJSON(e events.RequestCompleted, sizeStr, remote string) {
	entry := jsonEntry{
		Timestamp: e.Timestamp.Format(time.RFC3339),
		Method:    e.Method,
		Path:      e.Path,
		Status:    e.StatusCode,
		Size:      sizeStr,
		Duration:  e.Duration.Round(time.Microsecond).String(),
		Remote:    remote,
	}
	data, err := json.Marshal(entry)
	if err != nil {
		log.SetOutput(a.output)
		log.Printf("error marshaling log entry: %v", err)
		return
	}
	log.SetOutput(a.output)
	log.Println(string(data))
}
