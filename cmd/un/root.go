package main

import (
	"fmt"
	"time"

	"github.com/esnunes/un/pkg/cli"
	"github.com/esnunes/un/pkg/heredoc"
	"github.com/esnunes/un/pkg/ibge"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewRootCmd(log *logrus.Logger) *cobra.Command {
	dateFrom := cli.Date{Layout: cli.MonthLayout}
	dateTo := cli.Date{Layout: cli.MonthLayout}
	verbose := false

	cmd := &cobra.Command{
		Use:   "ipca [DATE_FROM] [DATE_TO]",
		Short: "Extended National Consumer Price Index (IPCA)",
		Long:  "The Extended National Consumer Price Index monthly variation",
		Example: heredoc.Doc(`
					# Returns the index for the current month
					ipca
					# Returns the index for the month 2022-08
					ipca 2022-08
					# Returns the index for the period between 2022-01 and 2022-07
					ipca 2022-01 2022-07
				`, 2),
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.RangeArgs(0, 2)(cmd, args); err != nil {
				return err
			}
			switch l := len(args); {
			case l == 0:
				dateFrom.Time = time.Now()
			case l > 0:
				if err := dateFrom.Set(args[0]); err != nil {
					return err
				}
				fallthrough
			case l > 1:
				if err := dateTo.Set(args[1]); err != nil {
					return err
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if verbose {
				log.SetLevel(logrus.InfoLevel)
			}

			c := ibge.Client{Log: log, Context: cmd.Context()}
			series, err := c.IPCA(dateFrom.Time, dateTo.Time)
			if err != nil {
				return err
			}

			w := cmd.OutOrStdout()
			for _, s := range series {
				fmt.Fprintf(w, "%s\t%7.2f\n", s.Date.Format(ibge.YearMonthLayout), s.Value)
			}

			return nil
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", verbose, "make the app talkative")
	return cmd
}
