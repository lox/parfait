package cmd

import (
	"fmt"

	"github.com/lox/parfait/api"
	"gopkg.in/alecthomas/kingpin.v2"
)

func ConfigureListStackOutputs(app *kingpin.Application, svc api.Services) {
	var stackName string

	cmd := app.Command("list-stack-outputs", "List all cloudformation stacks")
	cmd.Alias("outputs")

	cmd.Arg("name", "The name of the cloudformation stack").
		Required().
		StringVar(&stackName)

	cmd.Action(func(c *kingpin.ParseContext) error {
		outputs, err := api.StackOutputs(svc.Cloudformation, stackName)
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
