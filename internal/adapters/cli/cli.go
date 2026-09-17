// Package cli is a thin adapter that parses command-line flags into
// an application.Config. It has no business logic — validation and
// handler construction happen in the application layer.
package cli

import (
	"flag"
	"fmt"
	"io"
	"strconv"

	"github.com/aognio/webknife/internal/build"
)

// Config mirrors application.Config but is only produced by parsing.
// See internal/application for the canonical definition.
type Config struct {
	Command string

	ListenAddr string
	Root       string
	Upstream   string
	Auth       string
	TLSCert    string
	TLSKey     string
	LogFormat  string
	Verbose    bool

	EchoMaxBody int64

	RespondStatusCode  int
	RespondBody        string
	RespondBodyFile    string
	RespondContentType string

	RedirectTarget       string
	RedirectStatusCode   int
	RedirectPreservePath bool

	SetRequestHeaders     []string
	AddRequestHeaders     []string
	RemoveRequestHeaders  []string
	SetResponseHeaders    []string
	AddResponseHeaders    []string
	RemoveResponseHeaders []string
}

// Parse parses command-line args and returns the command name and config.
// It only handles flag parsing and basic required-field checks.
func Parse(args []string, stdout, stderr io.Writer) (string, *Config, error) {
	if len(args) < 2 {
		printUsage(stdout)
		return "", nil, fmt.Errorf("no command specified")
	}

	cmd := args[1]
	switch cmd {
	case "serve":
		return parseServe(args[2:], stdout)
	case "proxy":
		return parseProxy(args[2:], stdout)
	case "echo":
		return parseEcho(args[2:], stdout)
	case "respond":
		return parseRespond(args[2:], stdout)
	case "redirect":
		return parseRedirect(args[2:], stdout)
	case "version":
		fmt.Fprintf(stdout, "webknife %s (commit=%s, built=%s, go=%s)\n",
			build.Version, build.Commit, build.Date, build.GoVersion())
		return "", nil, nil
	case "help":
		printUsage(stdout)
		return "", nil, nil
	default:
		fmt.Fprintf(stderr, "unknown command: %s\n", cmd)
		printUsage(stderr)
		return "", nil, fmt.Errorf("unknown command: %s", cmd)
	}
}

func addCommonServerFlags(fs *flag.FlagSet, cfg *Config) {
	fs.StringVar(&cfg.ListenAddr, "listen", ":8080", "listen address")
	fs.StringVar(&cfg.Auth, "auth", "", "basic auth credentials (user:password)")
	fs.StringVar(&cfg.TLSCert, "tls-cert", "", "TLS certificate file")
	fs.StringVar(&cfg.TLSKey, "tls-key", "", "TLS private key file")
	fs.StringVar(&cfg.LogFormat, "log-format", "text", "log format (text|json)")
	fs.BoolVar(&cfg.Verbose, "v", false, "verbose output")
	fs.BoolVar(&cfg.Verbose, "verbose", false, "verbose output")
}

func addHeaderFlags(fs *flag.FlagSet, cfg *Config) {
	fs.Var(&multiStringValue{&cfg.SetRequestHeaders}, "set-request-header", "set request header (Name: Value)")
	fs.Var(&multiStringValue{&cfg.AddRequestHeaders}, "add-request-header", "add request header (Name: Value)")
	fs.Var(&multiStringValue{&cfg.RemoveRequestHeaders}, "remove-request-header", "remove request header (Name)")
	fs.Var(&multiStringValue{&cfg.SetResponseHeaders}, "set-response-header", "set response header (Name: Value)")
	fs.Var(&multiStringValue{&cfg.AddResponseHeaders}, "add-response-header", "add response header (Name: Value)")
	fs.Var(&multiStringValue{&cfg.RemoveResponseHeaders}, "remove-response-header", "remove response header (Name)")
}

func parseServe(args []string, stdout io.Writer) (string, *Config, error) {
	cfg := &Config{}
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	fs.SetOutput(stdout)
	addCommonServerFlags(fs, cfg)
	addHeaderFlags(fs, cfg)
	fs.StringVar(&cfg.Root, "root", ".", "document root directory")
	if err := fs.Parse(args); err != nil {
		return "", nil, err
	}
	return "serve", cfg, nil
}

func parseProxy(args []string, stdout io.Writer) (string, *Config, error) {
	cfg := &Config{}
	fs := flag.NewFlagSet("proxy", flag.ContinueOnError)
	fs.SetOutput(stdout)
	addCommonServerFlags(fs, cfg)
	addHeaderFlags(fs, cfg)
	fs.StringVar(&cfg.Upstream, "upstream", "", "upstream URL")
	if err := fs.Parse(args); err != nil {
		return "", nil, err
	}
	if cfg.Upstream == "" {
		return "", nil, fmt.Errorf("--upstream is required")
	}
	return "proxy", cfg, nil
}

func parseEcho(args []string, stdout io.Writer) (string, *Config, error) {
	cfg := &Config{}
	fs := flag.NewFlagSet("echo", flag.ContinueOnError)
	fs.SetOutput(stdout)
	addCommonServerFlags(fs, cfg)
	addHeaderFlags(fs, cfg)
	fs.Int64Var(&cfg.EchoMaxBody, "max-body", 1<<20, "max body size for inspection (bytes)")
	if err := fs.Parse(args); err != nil {
		return "", nil, err
	}
	return "echo", cfg, nil
}

func parseRespond(args []string, stdout io.Writer) (string, *Config, error) {
	cfg := &Config{}
	fs := flag.NewFlagSet("respond", flag.ContinueOnError)
	fs.SetOutput(stdout)
	addCommonServerFlags(fs, cfg)
	addHeaderFlags(fs, cfg)
	fs.IntVar(&cfg.RespondStatusCode, "status", 200, "HTTP status code")
	fs.StringVar(&cfg.RespondBody, "body", "", "response body")
	fs.StringVar(&cfg.RespondBodyFile, "body-file", "", "response body from file")
	fs.StringVar(&cfg.RespondContentType, "content-type", "", "Content-Type header")
	if err := fs.Parse(args); err != nil {
		return "", nil, err
	}
	return "respond", cfg, nil
}

func parseRedirect(args []string, stdout io.Writer) (string, *Config, error) {
	cfg := &Config{}
	fs := flag.NewFlagSet("redirect", flag.ContinueOnError)
	fs.SetOutput(stdout)
	addCommonServerFlags(fs, cfg)
	addHeaderFlags(fs, cfg)
	fs.StringVar(&cfg.RedirectTarget, "to", "", "redirect target URL")
	fs.IntVar(&cfg.RedirectStatusCode, "status", 302, "redirect status code (301,302,303,307,308)")
	fs.BoolVar(&cfg.RedirectPreservePath, "preserve-path", false, "preserve original path in redirect")
	if err := fs.Parse(args); err != nil {
		return "", nil, err
	}
	if cfg.RedirectTarget == "" {
		return "", nil, fmt.Errorf("--to is required")
	}
	return "redirect", cfg, nil
}

func printUsage(w io.Writer) {
	fmt.Fprintf(w, "webknife %s\n\n", build.Version)
	fmt.Fprintln(w, `Usage:
  webknife <command> [options]

Commands:
  serve     Serve static files
  proxy     Reverse proxy to an upstream
  echo      Inspect incoming requests
  respond   Return an arbitrary HTTP response
  redirect  Redirect to a URL
  version   Show version
  help      Show this help

Common options:
  --listen          Listen address (default ":8080")
  --auth            Basic auth credentials (user:password)
  --tls-cert        TLS certificate file
  --tls-key         TLS private key file
  --log-format      Log format: text or json (default "text")
  -v, --verbose     Verbose output

Header manipulation:
  --set-request-header 'Name: Value'
  --add-request-header 'Name: Value'
  --remove-request-header Name
  --set-response-header 'Name: Value'
  --add-response-header 'Name: Value'
  --remove-response-header Name

Examples:
  webknife serve --listen :8080 --root ./public
  webknife proxy --listen :8080 --upstream http://localhost:3000
  webknife echo --listen :8080
  webknife respond --listen :8080 --status 503 --body 'Service unavailable'
  webknife redirect --listen :8080 --to https://example.com
  webknife proxy --listen :8080 --upstream http://localhost:3000 \
    --set-request-header 'X-Debug: true' \
    --remove-response-header Server`)
}

type multiStringValue struct {
	target *[]string
}

func (m *multiStringValue) String() string {
	if m.target == nil || len(*m.target) == 0 {
		return ""
	}
	return strconv.Quote((*m.target)[len(*m.target)-1])
}

func (m *multiStringValue) Set(value string) error {
	*m.target = append(*m.target, value)
	return nil
}
