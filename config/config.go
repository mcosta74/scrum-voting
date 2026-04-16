package config

import (
	"flag"
	"log/slog"
	"os"
	"time"
)

type LogOptions struct {
	Level slog.Level `json:"level,omitempty"`
	UTC   bool       `json:"utc,omitempty"`
	JSON  bool       `json:"json,omitempty"`
}

type HTTPOptions struct {
	Addr    string        `json:"addr,omitempty"`
	Timeout time.Duration `json:"timeout,omitempty"`
}
type Options struct {
	Log  LogOptions  `json:"log,omitempty"`
	HTTP HTTPOptions `json:"http,omitempty"`
}

var defaultOptions = Options{
	Log: LogOptions{
		Level: slog.LevelInfo,
		UTC:   true,
		JSON:  true,
	},
	HTTP: HTTPOptions{
		Addr:    ":8080",
		Timeout: 10 * time.Second,
	},
}

func LoadOptions(appName string) (*Options, error) {
	opts := &Options{}
	*opts = defaultOptions

	fs := flag.NewFlagSet(appName, flag.ExitOnError)

	fs.TextVar(&opts.Log.Level, "log.level", opts.Log.Level, "The application's log level.")
	fs.BoolVar(&opts.Log.UTC, "log.utc", opts.Log.UTC, "Whether to use UTC for the timestamp of log messages.")
	fs.BoolVar(&opts.Log.JSON, "log.json", opts.Log.JSON, "Whether to use JSON format for the log messages.")

	fs.StringVar(&opts.HTTP.Addr, "http.addr", opts.HTTP.Addr, "Listen address of the HTTP server")
	fs.DurationVar(&opts.HTTP.Timeout, "http.timeout", opts.HTTP.Timeout, "Graceful shutdown timeout")

	_ = fs.Parse(os.Args[1:])
	return opts, nil
}
