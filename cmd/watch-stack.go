package cmd

import (
	"log"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/lox/parfait/api"
	"gopkg.in/alecthomas/kingpin.v2"
)

func ConfigureWatchStack(app *kingpin.Application, svc api.Services) {
	var stackName string

	cmd := app.Command("watch-stack", "Watch a Cloudformation stack until in a terminal state")
	cmd.Alias("w")

	cmd.Arg("name", "The name of the cloudformation stack to watch").
		StringVar(&stackName)

	cmd.Action(func(c *kingpin.ParseContext) error {
		return watchStack(svc, stackName)
	})
}

func watchStack(svc api.Services, stackName string) error {
	err := api.PollUntilCreated(svc.Cloudformation, stackName, func(event *cloudformation.StackEvent) {
		log.Printf("%s\n", api.FormatStackEvent(event))
	})
	if err != nil {
		return err
	}

	outputs, err := api.StackOutputs(svc.Cloudformation, stackName)
	if err != nil {
		return err
	}

	for k, v := range outputs {
		log.Printf("Stack Output: %s = %s", k, v)
	}

	return nil
}
