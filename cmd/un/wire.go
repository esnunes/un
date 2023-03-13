//go:build wireinject

package main

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/google/wire"
	"github.com/sirupsen/logrus"

	"github.com/esnunes/un/cmd/un/finance"
	"github.com/esnunes/un/cmd/un/ipca"
	"github.com/esnunes/un/cmd/un/openai"
	"github.com/esnunes/un/pkg/ibge"
	gpt "github.com/sashabaranov/go-gpt3"
)

func InitApp(ctx context.Context, gptClient *gpt.Client) *App {
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

		// cmd/un/openai
		openai.NewRootCmd,
		openai.NewConciseCmd,
		wire.Struct(new(openai.ConciseOptions), "*"),

		// cmd/un/finance
		finance.NewRootCmd,
		finance.NewRentalTaxCmd,
		wire.Struct(new(finance.RentalTaxOptions), "*"),

		// pkg/ibge
		wire.Struct(new(ibge.Client), "Log", "HTTP"),
		wire.Struct(new(ibge.IPCA), "*"),
	)
	return &App{}
}
