package cmd

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/lox/parfait/cmd/args"
	"github.com/lox/parfait/stacks"
	"gopkg.in/alecthomas/kingpin.v2"
)

func ConfigureCreateStack(app *kingpin.Application, sess client.ConfigProvider) {
	var stackName string
	var params []string
	var disableRollback bool

	cmd := app.Command("create-stack", "Create a cloudformation stack")
	cmd.Alias("create")

	tpl := args.TemplateSource(cmd.Flag("tpl", "Either a file path or url to a cloudformation template").
		Short('t'))

	cmd.Flag("no-rollback", "Disable stack rollback on failure").
		BoolVar(&disableRollback)

	cmd.Arg("stack-name", "The name of the cloudformation stack").
		StringVar(&stackName)

	cmd.Arg("params", "Parameters to the stack in Key=Val form").
		StringsVar(&params)

	cmd.Action(func(c *kingpin.ParseContext) error {
		params, err := parseStackParams(params)
		if err != nil {
			return err
		}

		ctx := stacks.CreateStackContext{
			Params:          params,
			Body:            tpl.String(),
			DisableRollback: disableRollback,
		}

		cfn := cloudformation.New(sess)

		if err = stacks.Create(cfn, stackName, ctx); err != nil {
			return err
		}

		return stacks.Watch(cfn, stackName, func(event *cloudformation.StackEvent) {
			fmt.Printf("%s\n", stacks.FormatStackEvent(event))
		})
	})
}

func parseStackParams(rawParams []string) (map[string]string, error) {
	params := map[string]string{}
	for _, arg := range rawParams {
		parts := strings.Split(arg, "=")
		if len(parts) != 2 {
			return params, fmt.Errorf("Failed to parse parameter %q", arg)
		}
		params[parts[0]] = parts[1]
	}
	return params, nil
}
