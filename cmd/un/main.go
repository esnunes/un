package main

import (
	"context"
	"os"

	"github.com/esnunes/un/pkg/shutdown"
	gpt "github.com/sashabaranov/go-gpt3"
	"github.com/sirupsen/logrus"
)

type App struct {
	Log     *logrus.Logger
	RootCmd *RootCmd
}

func main() {
	ctx, cancel := shutdown.Context(context.Background())
	defer cancel()

	app := InitApp(ctx, gpt.NewClient(os.Getenv("OPENAI_KEY")))
	app.Log.SetLevel(logrus.WarnLevel)

	if err := app.RootCmd.Execute(); err != nil {
		app.Log.Errorf("failed to execute command: %v", err)
		if ctx.Err() != nil {
			app.Log.Errorf("interrupted: %v\n", ctx.Err())
			if sig := shutdown.SignalFromContext(ctx); sig != nil {
				app.Log.Errorf("signal: %v\n", *sig)
			}
		}
		os.Exit(1)
	}
}
