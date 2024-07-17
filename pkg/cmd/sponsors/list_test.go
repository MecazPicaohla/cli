package sponsors_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/v2/pkg/cmd/sponsors"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/google/shlex"
	"github.com/stretchr/testify/require"
)

func TestNewCmdSponsors(t *testing.T) {
	tests := []struct {
		name            string
		args            string
		expectedErr     error
		expectedOptions sponsors.ListOptions
	}{
		{
			name:        "when no arguments provided, returns a useful error",
			args:        "",
			expectedErr: cmdutil.FlagErrorf("must specify a user"),
		},
		{
			name: "org",
			args: "testusername",
			expectedOptions: sponsors.ListOptions{
				User: "testusername",
			},
		},
		{
			name:        "when too many arguments provided, returns a useful error",
			args:        "foo bar",
			expectedErr: cmdutil.FlagErrorf("too many arguments"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &cmdutil.Factory{
				HttpClient: func() (*http.Client, error) {
					return nil, nil
				},
			}

			argv, err := shlex.Split(tt.args)
			require.NoError(t, err)

			var optsSpy sponsors.ListOptions
			cmd := sponsors.NewCmdList(f, func(opts sponsors.ListOptions) error {
				optsSpy = opts
				return nil
			})
			cmd.SetArgs(argv)

			_, err = cmd.ExecuteC()
			require.Equal(t, tt.expectedErr, err)
			require.Equal(t, tt.expectedOptions.User, optsSpy.User)
		})
	}
}

type FakeSponsorLister struct {
	StubbedSponsors map[sponsors.User][]sponsors.Sponsor
	StubbedErr      error
}

func (s FakeSponsorLister) ListSponsors(user sponsors.User) ([]sponsors.Sponsor, error) {
	if s.StubbedErr != nil {
		return nil, s.StubbedErr
	}
	return s.StubbedSponsors[user], nil
}

func TestCmdSponsorsListPrintsTableToTTY(t *testing.T) {
	// Given our sponsor lister returns successfully
	ios, _, stdout, _ := iostreams.Test()
	ios.SetStdoutTTY(true)
	ios.SetStdinTTY(true)
	ios.SetStderrTTY(true)
	listOptions := sponsors.ListOptions{
		IO: ios,
		SponsorLister: FakeSponsorLister{
			StubbedSponsors: map[sponsors.User][]sponsors.Sponsor{
				"testusername": {"sponsor1", "sponsor2"},
			},
		},
		User: "testusername",
	}

	// When I run the list command
	err := sponsors.ListRun(listOptions)

	// Then it is successful
	require.NoError(t, err)

	// And it pretty prints a table containing the sponsor names to our TTY
	expectedOutput := heredoc.Doc(`
	USER
	sponsor1
	sponsor2
	`)
	require.Equal(t, expectedOutput, stdout.String())
}

func TestCmdSponsorsListSponsorListingError(t *testing.T) {
	// Given our sponsor lister returns with an error
	ios, _, _, _ := iostreams.Test()
	ios.SetStdoutTTY(true)
	ios.SetStdinTTY(true)
	ios.SetStderrTTY(true)
	listOptions := sponsors.ListOptions{
		IO: ios,
		SponsorLister: FakeSponsorLister{
			StubbedErr: errors.New("expected test error"),
		},
		User: "testusername",
	}

	// When I run the list command
	err := sponsors.ListRun(listOptions)

	// Then it returns an informational error
	require.ErrorContains(t, err, "sponsor list: expected test error")
}
