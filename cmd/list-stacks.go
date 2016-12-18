package cmd

import (
	"fmt"

	"github.com/lox/parfait/api"
	"gopkg.in/alecthomas/kingpin.v2"
)

func ConfigureListStacks(app *kingpin.Application, svc api.Services) {
	cmd := app.Command("list-stacks", "List all cloudformation stacks")
	cmd.Alias("list")
	cmd.Alias("ls")

	cmd.Action(func(c *kingpin.ParseContext) error {
		stacks, err := api.FindAllActiveStacks(svc.Cloudformation)
		if err != nil {
			return err
		}

		fmt.Printf("%-60s %-40s %-20s\n", "NAME", "STATUS", "LAST UPDATED")
		for _, stack := range stacks {
			t := *stack.CreationTime
			if stack.LastUpdatedTime != nil {
				t = *stack.LastUpdatedTime
			}

			fmt.Printf("%-60s %-40s %-20s\n",
				*stack.StackName,
				*stack.StackStatus,
				t.String(),
			)
		}
		return nil
	})
}
