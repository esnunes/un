package main

import (
	"github.com/esnunes/un/cmd/un/ipca"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type RootCmd struct{ *cobra.Command }

func NewRootCmd(log *logrus.Logger, ipca *ipca.RootCmd) *RootCmd {
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

	return &RootCmd{cmd}
}
