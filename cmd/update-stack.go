package cmd

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/lox/parfait/cmd/args"
	"github.com/lox/parfait/stacks"
	"gopkg.in/alecthomas/kingpin.v2"
)

func ConfigureUpdateStack(app *kingpin.Application, sess client.ConfigProvider) {
	var stackName string
	var params []string

	cmd := app.Command("update-stack", "Update a cloudformation stack")
	cmd.Alias("update")

	tpl := args.TemplateSource(cmd.Flag("tpl", "Either a file path or url to a cloudformation template").
		Short('t'))

	cmd.Arg("name", "The name of the cloudformation stack").
		Required().
		StringVar(&stackName)

	cmd.Arg("params", "Parameters to the stack in Key=Val form").
		StringsVar(&params)

	cmd.Action(func(c *kingpin.ParseContext) error {
		params, err := parseStackParams(params)
		if err != nil {
			return err
		}

		ctx := stacks.CreateStackContext{
			Params: params,
			Body:   tpl.String(),
		}

		t := time.Now()
		svc := cloudformation.New(sess)

		if err = stacks.Create(svc, stackName, ctx); err != nil {
			return err
		}

		return stacks.Watch(svc, stackName, func(event *cloudformation.StackEvent) {
			if event.Timestamp.After(t) {
				fmt.Printf("%s\n", stacks.FormatStackEvent(event))
			}
		})
	})
}
