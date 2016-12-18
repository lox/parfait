package main

import (
	"os"

	"github.com/lox/parfait/api"
	"github.com/lox/parfait/cmd"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	Version string
)

func main() {
	run(os.Args[1:], os.Exit)
}

func run(args []string, exit func(code int)) {
	app := kingpin.New("parfait",
		`A tool for creating and monitoring AWS CloudFormation stacks`)

	app.Version(Version)
	app.Writer(os.Stdout)
	app.DefaultEnvars()
	app.Terminate(exit)

	cmd.ConfigureWatchStack(app, api.DefaultServices)
	cmd.ConfigureListStacks(app, api.DefaultServices)
	cmd.ConfigureListStackOutputs(app, api.DefaultServices)
	cmd.ConfigureCreateStack(app, api.DefaultServices)
	cmd.ConfigureUpdateStack(app, api.DefaultServices)
	cmd.ConfigureFollowLogs(app, api.DefaultServices)

	kingpin.MustParse(app.Parse(args))
}
