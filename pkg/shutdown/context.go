package shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type ctxType string

const (
	signalKey ctxType = "signal"
)

func Context(ctx context.Context) (context.Context, context.CancelFunc) {
	var sigValue os.Signal

	ctx = context.WithValue(ctx, signalKey, &sigValue)
	ctx, cancel := context.WithCancel(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		select {
		case sig := <-c:
			sigValue = sig
			cancel()
		case <-ctx.Done():
		}
	}()

	return ctx, cancel
}

func SignalFromContext(ctx context.Context) *os.Signal {
	return ctx.Value(signalKey).(*os.Signal)
}
