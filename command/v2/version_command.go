package v2

import (
	"os"

	"code.cloudfoundry.org/cli/cf/cmd"
	"code.cloudfoundry.org/cli/command"
)

type VersionCommand struct {
	usage interface{} `usage:"CF_NAME version\n\n   'cf -v' and 'cf --version' are also accepted."`
}

func (_ VersionCommand) Setup(config command.Config, ui command.UI) error {
	return nil
}

func (_ VersionCommand) Execute(args []string) error {
	cmd.Main(os.Getenv("CF_TRACE"), os.Args)
	return nil
}
