package finance

import (
	"context"
	"fmt"
	"io"
	"math/big"

	"github.com/esnunes/un/pkg/heredoc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type RentalTaxCmd struct{ *cobra.Command }

func NewRentalTaxCmd(log *logrus.Logger, opts *RentalTaxOptions) *RentalTaxCmd {
	cmd := &cobra.Command{
		Use:  "rental-tax VALUE",
		Long: "Calculate the rental tax based on the given value",
		Example: heredoc.Doc(`
					# Returns the rental tax for the value 4530
					rental-tax 4530
				`, 2),
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.Args = func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			value, _, err := big.NewFloat(0).Parse(args[0], 10)
			if err != nil {
				return err
			}
			opts.Value = value
		}
		return cobra.MinimumNArgs(1)(cmd, args)
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return RentalTax(opts)
	}

	return &RentalTaxCmd{cmd}
}

type RentalTaxOptions struct {
	Context context.Context
	Out     io.Writer
	Value   *big.Float `wire:"-"`
}

type TaxRule struct {
	Threshold *big.Float
	Rate      *big.Float
	Base      *big.Float
}

var (
	rules = []TaxRule{
		{Threshold: big.NewFloat(4664.68), Rate: big.NewFloat(0.275), Base: big.NewFloat(413.43)},
		{Threshold: big.NewFloat(3751.05), Rate: big.NewFloat(0.225), Base: big.NewFloat(207.86)},
		{Threshold: big.NewFloat(2826.65), Rate: big.NewFloat(0.15), Base: big.NewFloat(69.20)},
		{Threshold: big.NewFloat(1903.98), Rate: big.NewFloat(0.075), Base: big.NewFloat(0)},
	}
)

func RentalTax(o *RentalTaxOptions) error {
	value := big.NewFloat(0.0)
	for _, r := range rules {
		if o.Value.Cmp(r.Threshold) > 0 {
			value = value.Sub(o.Value, r.Threshold)
			value = value.Mul(value, r.Rate)
			value = value.Add(value, r.Base)
			break
		}
	}
	fmt.Fprintf(o.Out, "%v\n", value.Text('f', 2))
	return nil
}
