package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"syscall"

	"github.com/mcosta74/scrum-voting/config"
	"github.com/mcosta74/scrum-voting/handlers"
	"github.com/oklog/run"
)

//go:embed webapp/dist/webapp/browser
var staticFS embed.FS

func main() {
	opts, err := config.LoadOptions("scrum-voting")
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	var (
		ctx         = context.Background()
		logger      = newLogger(opts.Log)
		httpHandler = handlers.NewHandler(staticFS)
	)

	logger.Info("service started")
	defer logger.Info("service stopped")

	entries, err := staticFS.ReadDir(".")
	if err == nil {
		for _, entry := range entries {
			fmt.Printf("%+v\n", entry)
		}
	}
	var g run.Group
	{
		g.Add(run.SignalHandler(ctx, syscall.SIGINT, syscall.SIGTERM))
	}
	{
		httpServer := &http.Server{
			Addr:    opts.HTTP.Addr,
			Handler: httpHandler,
		}

		g.Add(func() error {
			logger.Info("starting HTTP server...", "addr", httpServer.Addr)
			return httpServer.ListenAndServe()
		}, func(err error) {
			ctx, cancel := context.WithTimeout(ctx, opts.HTTP.Timeout)
			defer cancel()

			if err := httpServer.Shutdown(ctx); err != nil {
				if errors.Is(err, context.Canceled) {
					logger.Error("failure gracefully shutting down the HTTP server", "reason", err)
				}
			}
		})
	}
	logger.Info("service shutdown", "reason", g.Run())
}

func newLogger(opts config.LogOptions) *slog.Logger {
	logOpts := &slog.HandlerOptions{
		Level:     opts.Level,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && opts.UTC {
				a.Value = slog.TimeValue(a.Value.Time().UTC())
			}

			if a.Key == slog.SourceKey {
				source := a.Value.Any().(*slog.Source)
				source.File = filepath.Base(source.File)
			}
			return a
		},
	}

	var h slog.Handler
	if opts.JSON {
		h = slog.NewJSONHandler(os.Stdout, logOpts)
	} else {
		h = slog.NewTextHandler(os.Stdout, logOpts)
	}
	return slog.New(h)
}
