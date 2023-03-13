package main

import (
	"github.com/esnunes/un/cmd/un/finance"
	"github.com/esnunes/un/cmd/un/ipca"
	"github.com/esnunes/un/cmd/un/openai"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type RootCmd struct{ *cobra.Command }

func NewRootCmd(
	log *logrus.Logger, ipca *ipca.RootCmd, openai *openai.RootCmd, finance *finance.RootCmd,
) *RootCmd {
	verbose := false

	cmd := &cobra.Command{
		Use:           "un",
		Long:          "The Uncle Nunes set of tools",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.NoArgs,
	}

	cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if verbose {
			log.SetLevel(logrus.InfoLevel)
		}
	}

	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", verbose, "make the app verbose")

	cmd.AddCommand(ipca.Command)
	cmd.AddCommand(openai.Command)
	cmd.AddCommand(finance.Command)

	return &RootCmd{cmd}
}
