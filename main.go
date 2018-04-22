package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/lox/parfait/cmd"
	"github.com/lox/parfait/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	run(os.Args[1:], os.Exit)
}

func run(args []string, exit func(code int)) {
	app := kingpin.New("parfait",
		`A tool for creating and monitoring AWS CloudFormation stacks`)

	app.Version(version.Version)
	app.Writer(os.Stdout)
	app.DefaultEnvars()
	app.Terminate(exit)

	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		log.Fatal(err)
	}

	cmd.ConfigureWatchStack(app, sess)
	cmd.ConfigureListStacks(app, sess)
	cmd.ConfigureListStackOutputs(app, sess)
	cmd.ConfigureCreateStack(app, sess)
	cmd.ConfigureUpdateStack(app, sess)
	cmd.ConfigureDeleteStack(app, sess)
	cmd.ConfigureFollowLogs(app, sess)

	kingpin.MustParse(app.Parse(args))
}
