package cmd

import (
	"log"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/lox/parfait/api"
	"gopkg.in/alecthomas/kingpin.v2"
)

func ConfigureWatchStack(app *kingpin.Application, svc api.Services) {
	var stackName string

	cmd := app.Command("watch", "Watch a supported AWS resource, either Cloudformation or Cloudwatch logs")
	cmd.Alias("w")

	cmd.Flag("stack", "The name of the cloudformation stack to watch").
		Required().
		StringVar(&stackName)

	cmd.Flag("stack", "The name of the cloudformation stack to watch").
		Required().
		StringVar(&stackName)

	cmd.Action(func(c *kingpin.ParseContext) error {
		return pollStack(svc, stackName)
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
