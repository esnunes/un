//go:build wireinject

package main

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/google/wire"
	"github.com/sirupsen/logrus"

	"github.com/esnunes/un/cmd/un/ipca"
	"github.com/esnunes/un/pkg/ibge"
)

func InitApp(ctx context.Context) *App {
	wire.Build(
		// general
		logrus.New,
		wire.Value(http.DefaultClient),
		wire.InterfaceValue(new(io.Writer), os.Stdout),

		// cmd/un
		NewRootCmd,
		wire.Struct(new(App), "*"),

		// cmd/un/ipca
		ipca.NewRootCmd,
		ipca.NewRateCmd,
		wire.Bind(new(ipca.RateRetriever), new(*ibge.IPCA)),
		wire.Struct(new(ipca.RateOptions), "*"),

		// pkg/ibge
		wire.Struct(new(ibge.Client), "Log", "HTTP"),
		wire.Struct(new(ibge.IPCA), "*"),
	)
	return &App{}
}
