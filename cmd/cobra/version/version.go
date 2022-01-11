package version

import (
	"github.com/rafaelcalleja/go-kit/cmd/helper"
	"github.com/rafaelcalleja/go-kit/cmd/termcolor"
	"github.com/rafaelcalleja/go-kit/logger"
	"github.com/spf13/cobra"
)

// Build information. Populated at build-time.
var (
	Version   string
	Revision  string
	Branch    string
	BuildUser string
	BuildDate string
	GoVersion string
)

const (

	// TestVersion used in test cases for the current version if no
	// version can be found - such as if the version property is not properly
	// included in the go test flags
	TestVersion = "1.0.0-SNAPSHOT"
)

// ShowOptions the options for viewing running PRs
type Options struct {
	Verbose bool
	Quiet   bool
}

// NewCmdVersion creates a command object for the "version" command
func NewCmdVersion(helper helper.ErrorHelper, logger logger.Logger, termColor termcolor.TermColor) *cobra.Command {
	o := &Options{}

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Displays the version of this command",
		Run: func(cmd *cobra.Command, args []string) {
			err := o.Run(logger, termColor)
			helper.CheckErr(err)
		},
	}

	cmd.Flags().BoolVarP(&o.Quiet, "quiet", "q", false, "uses the quiet format of just outputting the version number only")

	return cmd
}

// Run implements the command
func (o *Options) Run(logger logger.Logger, termColor termcolor.TermColor) error {
	v := GetVersion()
	if o.Quiet {
		logger.Infof(v)
		return nil
	}
	logger.Infof("version: %s", termColor.ColorInfo(v))
	return nil
}

func GetVersion() string {
	if Version != "" {
		return Version
	}
	return TestVersion
}
