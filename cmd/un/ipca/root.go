package ipca

import (
	"github.com/spf13/cobra"
)

type RootCmd struct{ *cobra.Command }

func NewRootCmd(rate *RateCmd) *RootCmd {
	cmd := &cobra.Command{
		Use:           "ipca",
		Long:          "Extended National Consumer Price Index (IPCA)",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.NoArgs,
	}

	cmd.AddCommand(rate.Command)

	return &RootCmd{cmd}
}
