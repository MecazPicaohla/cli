package sponsors

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/v2/api"
	"github.com/cli/cli/v2/internal/tableprinter"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
)

type User string

type ListOptions struct {
	IO *iostreams.IOStreams

	SponsorLister SponsorLister

	User User
}

func NewCmdList(f *cmdutil.Factory, runF func(ListOptions) error) *cobra.Command {
	opts := ListOptions{
		IO: f.IOStreams,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List sponsors for a user",
		Example: heredoc.Doc(`
			List sponsors of a user
			$ gh sponsors list <user>
		`),
		Aliases: []string{"ls"},
		Args:    cmdutil.ExactArgs(1, "must specify a user"),
		RunE: func(cmd *cobra.Command, args []string) error {
			httpClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			opts.SponsorLister = GQLSponsorClient{
				Hostname:  "github.com",
				APIClient: api.NewClientFromHTTP(httpClient),
			}

			opts.User = User(args[0])

			if runF != nil {
				return runF(opts)
			}
			return ListRun(opts)
		},
	}

	return cmd
}

type Sponsor string

type SponsorLister interface {
	ListSponsors(user User) ([]Sponsor, error)
}

func ListRun(opts ListOptions) error {
	sponsors, err := opts.SponsorLister.ListSponsors(opts.User)
	if err != nil {
		return fmt.Errorf("sponsor list: %v", err)
	}

	tp := tableprinter.New(opts.IO, tableprinter.WithHeader("USER"))
	for _, sponsor := range sponsors {
		tp.AddField(string(sponsor))
		tp.EndRow()
	}

	return tp.Render()
}
