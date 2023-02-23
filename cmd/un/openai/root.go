package openai

import (
	"github.com/spf13/cobra"
)

type RootCmd struct{ *cobra.Command }

func NewRootCmd(concise *ConciseCmd) *RootCmd {
	cmd := &cobra.Command{
		Use:           "openai",
		Long:          "OpenAI API",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.NoArgs,
	}

	cmd.AddCommand(concise.Command)

	return &RootCmd{cmd}
}
