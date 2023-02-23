package openai

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/esnunes/un/pkg/heredoc"
	gpt "github.com/sashabaranov/go-gpt3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type ConciseCmd struct{ *cobra.Command }

func NewConciseCmd(log *logrus.Logger, opts *ConciseOptions) *ConciseCmd {
	cmd := &cobra.Command{
		Use:  "concise TEXT",
		Long: "Make the TEXT more concise",
		Example: heredoc.Doc(`
					# Returns the concise version of the given text
					openai concise a very long text
				`, 2),
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.Args = func(cmd *cobra.Command, args []string) error {
		opts.Text = strings.Join(args, " ")
		return cobra.MinimumNArgs(1)(cmd, args)
	}

	cmd.Flags().IntVarP(&opts.MaxTokens, "max-tokens", "t", 60, "maximum number of words/tokens returned")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return Concise(opts)
	}

	return &ConciseCmd{cmd}
}

// type RateRetriever interface {
// 	Rate(ctx context.Context, from, to time.Time) ([]ibge.Serie, error)
// }

type ConciseOptions struct {
	Context   context.Context
	Out       io.Writer
	Client    *gpt.Client
	MaxTokens int    `wire:"-"`
	Text      string `wire:"-"`
}

func Concise(o *ConciseOptions) error {
	r, err := o.Client.CreateCompletion(o.Context, gpt.CompletionRequest{
		Model:     gpt.GPT3TextDavinci002,
		Prompt:    "Make the following text more concise:\n" + o.Text,
		MaxTokens: o.MaxTokens,
	})
	if err != nil {
		return err
	}

	if len(r.Choices) == 0 {
		return errors.New("it could not make the text more concise")
	}

	fmt.Fprintln(o.Out, strings.Trim(r.Choices[0].Text, "\n"))

	return nil
}
