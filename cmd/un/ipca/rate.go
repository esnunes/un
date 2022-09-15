package ipca

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/esnunes/un/pkg/cli"
	"github.com/esnunes/un/pkg/heredoc"
	"github.com/esnunes/un/pkg/ibge"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type RateCmd struct{ *cobra.Command }

func NewRateCmd(log *logrus.Logger, opts *RateOptions) *RateCmd {
	cmd := &cobra.Command{
		Use:  "rate [DATE_FROM] [DATE_TO]",
		Long: "The Extended National Consumer Price Index rate monthly variation",
		Example: heredoc.Doc(`
					# Returns the index for the current month
					ipca rate
					# Returns the index for the month 2022-08
					ipca rate 2022-08
					# Returns the index for the period between 2022-01 and 2022-07
					ipca rate 2022-01 2022-07
				`, 2),
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.Args = func(cmd *cobra.Command, args []string) error {
		if err := cobra.RangeArgs(0, 2)(cmd, args); err != nil {
			return err
		}
		switch l := len(args); {
		case l == 0:
			opts.DateFrom = time.Now()
		case l > 1:
			date := cli.Date{Layout: cli.MonthLayout}
			if err := date.Set(args[1]); err != nil {
				return err
			}
			opts.DateTo = date.Time
			fallthrough
		case l > 0:
			date := cli.Date{Layout: cli.MonthLayout}
			if err := date.Set(args[0]); err != nil {
				return err
			}
			opts.DateFrom = date.Time
		}
		return nil
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return Rate(opts)
	}

	return &RateCmd{cmd}
}

type RateRetriever interface {
	Rate(ctx context.Context, from, to time.Time) ([]ibge.Serie, error)
}

type RateOptions struct {
	Retriever        RateRetriever
	Context          context.Context
	Out              io.Writer
	DateFrom, DateTo time.Time `wire:"-"`
}

func Rate(o *RateOptions) error {
	series, err := o.Retriever.Rate(o.Context, o.DateFrom, o.DateTo)
	if err != nil {
		return err
	}

	for _, s := range series {
		fmt.Fprintf(o.Out, "%s\t%7.2f\n", s.Date.Format(ibge.YearMonthLayout), s.Value)
	}

	return nil
}
