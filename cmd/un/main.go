package main

import (
	"context"
	"os"

	"github.com/esnunes/un/pkg/shutdown"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetLevel(logrus.WarnLevel)

	ctx, cancel := shutdown.Context(context.Background())
	defer cancel()

	if err := NewRootCmd(log).ExecuteContext(ctx); err != nil {
		log.Errorf("failed to execute command: %v", err)
		if ctx.Err() != nil {
			log.Errorf("interrupted: %v\n", ctx.Err())
			if sig := shutdown.SignalFromContext(ctx); sig != nil {
				log.Errorf("signal: %v\n", *sig)
			}
		}
		os.Exit(1)
	}
}
