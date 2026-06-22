package main

import (
	"context"
	"errors"
	"fmt"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const v3RendererEbiten = "ebiten"

func normalizeV3Renderer(renderer string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(renderer)) {
	case "", v3RendererEbiten:
		return v3RendererEbiten, nil
	default:
		return "", fmt.Errorf("unsupported --renderer %q; expected ebiten", renderer)
	}
}

func newV3Context(duration time.Duration) (context.Context, func()) {
	signalCtx, stopSignals := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	if duration <= 0 {
		return signalCtx, stopSignals
	}
	ctx, cancel := context.WithTimeout(signalCtx, duration)
	return ctx, func() {
		cancel()
		stopSignals()
	}
}

func ignoreContextStop(err error) error {
	if isContextStop(err) {
		return nil
	}
	return err
}

func isContextStop(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}
