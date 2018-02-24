package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/fatih/color"
	"github.com/lox/parfait/stacks"
	"gopkg.in/alecthomas/kingpin.v2"
)

func ConfigureWatchStack(app *kingpin.Application, sess client.ConfigProvider) {
	var stackName string

	cmd := app.Command("watch-stack", "Watch a Cloudformation stack until in a terminal state")
	cmd.Alias("watch")
	cmd.Alias("w")

	cmd.Arg("name", "The name of the cloudformation stack to watch").
		StringVar(&stackName)

	cmd.Action(func(c *kingpin.ParseContext) error {
		err := stacks.Watch(cloudformation.New(sess), stackName, func(event *cloudformation.StackEvent) {
			fmt.Printf("%s\n", stacks.FormatStackEvent(event))
		})
		if err != nil {
			fmt.Printf("\n%v\n\n", color.RedString(err.Error()))
			os.Exit(1)
		}
		return nil
	})
}
