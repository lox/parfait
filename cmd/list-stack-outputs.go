package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/lox/parfait/stacks"
	"gopkg.in/alecthomas/kingpin.v2"
)

func ConfigureListStackOutputs(app *kingpin.Application, sess client.ConfigProvider) {
	var stackName string

	cmd := app.Command("list-stack-outputs", "List all cloudformation stacks")
	cmd.Alias("outputs")

	cmd.Arg("name", "The name of the cloudformation stack").
		Required().
		StringVar(&stackName)

	cmd.Action(func(c *kingpin.ParseContext) error {
		cfn := cloudformation.New(sess)

		outputs, err := stacks.Outputs(cfn, stackName)
		if err != nil {
			return err
		}

		fmt.Printf("%-20s %-80s\n", "KEY", "VALUE")
		for k, v := range outputs {
			fmt.Printf("%-20s %-80s\n", k, v)
		}
		return nil
	})
}
