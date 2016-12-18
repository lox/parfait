package cmd

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/lox/parfait/api"
	"gopkg.in/alecthomas/kingpin.v2"
)

func ConfigureDeleteStack(app *kingpin.Application, svc api.Services) {
	var stackName string

	cmd := app.Command("delete-stack", "Update a cloudformation stack")
	cmd.Alias("delete")
	cmd.Alias("del")
	cmd.Alias("remove")
	cmd.Alias("rm")

	cmd.Arg("name", "The name of the cloudformation stack").
		Required().
		StringVar(&stackName)

	cmd.Action(func(c *kingpin.ParseContext) error {
		t := time.Now()
		if err := api.DeleteStack(svc.Cloudformation, stackName); err != nil {
			return err
		}

		return api.PollUntilDeleted(svc.Cloudformation, stackName, func(event *cloudformation.StackEvent) {
			if event.Timestamp.After(t) {
				log.Printf("%s\n", api.FormatStackEvent(event))
			}
		})
	})
}
