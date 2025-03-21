package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/4925k/ardanlabs/software_design_with_kubernetes/service/app/services/sales/v1/handlers"
	v1 "github.com/4925k/ardanlabs/software_design_with_kubernetes/service/business/web/v1"
	"github.com/4925k/ardanlabs/software_design_with_kubernetes/service/business/web/v1/debug"
	"github.com/4925k/ardanlabs/software_design_with_kubernetes/service/foundation/logger"
	"github.com/4925k/ardanlabs/software_design_with_kubernetes/service/foundation/web"
	"github.com/ardanlabs/conf/v3"
)

var build = "develop"

func main() {
	// LOGGING
	var log *logger.Logger

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "********* SEND ALERT *********")
		},
	}

	traceIDFn := func(ctx context.Context) string {
		return web.GetTraceID(ctx)
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "SALES", traceIDFn, events)

	// --------------------------------------------------------------------------------------
	// RUN
	ctx := context.Background()

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "msg", err)
		return
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	// --------------------------------------------------------------------------------------
	// GOMAXPROCS

	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0), "build", build)
	expvar.NewString("build").Set(build)

	// --------------------------------------------------------------------------------------
	// CONFIGURATION

	cfg := struct {
		conf.Version
		Web struct {
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:20s,mask"` // noprint tage hides it from the help, mask hide the value
			APIHost         string        `conf:"default:0.0.0.0:3000"`
			DebugHost       string        `conf:"default:0.0.0.0:4000"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "4925k",
		},
	}

	const prefix = "SALES"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	// --------------------------------------------------------------------------------------
	// App Starting

	log.Info(ctx, "starting service", "version", build)
	defer log.Info(ctx, "shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}

	log.Info(ctx, "startup", "config", out)

	// --------------------------------------------------------------------------------------
	// Debug Service

	go func() {
		log.Info(ctx, "startup", "status", "debug v1 router started", "host", cfg.Web.DebugHost)

		// NOTE:
		// http.DefaultServeMux will host any paths that might be initialized by other packages
		// this leaves us with unknown paths
		// to address this we create a custom mux that only hosts the paths we care about
		if err := http.ListenAndServe(cfg.Web.DebugHost, debug.Mux()); err != nil {
			log.Error(ctx, "shutdown", "status", "debug v1 router closed", "host", cfg.Web.DebugHost, "msg", err)
		}
	}()

	// --------------------------------------------------------------------------------------
	// API SERVICE

	log.Info(ctx, "startup", "status", "initializing V1 API Support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	cfgMux := v1.APIMuxConfig{
		Build:    build,
		Shutdown: shutdown,
		Log:      log,
	}

	apiMux := v1.APIMux(cfgMux, handlers.Routes{})

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      apiMux,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// --------------------------------------------------------------------------------------
	// SHUTDOWN

	select {
	case err := <-serverErrors:
		return fmt.Errorf("listen and serve: %w", err)
	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
