package sponsors

import (
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdSponsors(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sponsors <command>",
		Short: "Manage sponsors",
		Long:  `Work with GitHub sponsors.`,
	}

	cmd.AddCommand(NewCmdList(f, nil))
	return cmd
}
