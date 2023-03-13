package finance

import (
	"github.com/spf13/cobra"
)

type RootCmd struct{ *cobra.Command }

func NewRootCmd(rentalTax *RentalTaxCmd) *RootCmd {
	cmd := &cobra.Command{
		Use:           "finance",
		Long:          "Finance related stuff",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.NoArgs,
	}

	cmd.AddCommand(rentalTax.Command)

	return &RootCmd{cmd}
}
