package cmd

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/lox/parfait/stacks"
	"github.com/lox/parfait/stacks/poller"
	"gopkg.in/alecthomas/kingpin.v2"
)

func ConfigureDeleteStack(app *kingpin.Application, sess client.ConfigProvider) {
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
		cfn := cloudformation.New(sess)

		t := time.Now()
		if err := stacks.Delete(cfn, stackName); err != nil {
			return err
		}

		return poller.UntilDeleted(cfn, stackName, func(event *cloudformation.StackEvent) {
			if event.Timestamp.After(t) {
				fmt.Printf("%s\n", stacks.FormatStackEvent(event))
			}
		})
	})
}
