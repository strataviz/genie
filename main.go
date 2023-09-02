package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"ctx.sh/genie/pkg/cmd"
	"ctx.sh/genie/pkg/config"
	"ctx.sh/strata"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGSTOP)
	defer cancel()

	// Logging
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	zapCfg := zap.Config{
		// TODO: make me configurable
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "console",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
	}
	zl := zap.Must(zapCfg.Build())
	defer zl.Sync() //nolint:errcheck

	// Metrics
	logger := zapr.NewLogger(zl)
	metrics := strata.New(strata.MetricsOpts{
		Logger:       logger,
		Prefix:       []string{"genie"},
		PanicOnError: true,
	})

	var obs sync.WaitGroup
	obs.Add(1)
	go func() {
		defer obs.Done()
		err := metrics.Start(strata.ServerOpts{
			Port: 9090,
		})
		if err != nil && err != http.ErrServerClosed {
			logger.Error(err, "metrics server start failed")
			os.Exit(1)
		}
	}()

	// Fix me
	cfg, err := config.Load([]string{"./genie.d"})
	if err != nil {
		logger.Error(err, "unable to load configuration")
		os.Exit(1)
	}

	root := cmd.NewRoot(&cmd.GlobalOpts{
		Logger:      logger,
		Metrics:     metrics,
		BaseContext: ctx,
		CancelFunc:  cancel,
		Config:      cfg,
	})
	root.Execute()

	// Shut down observability tools.
	metrics.Stop()
	obs.Wait()
}
