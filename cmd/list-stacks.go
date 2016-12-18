package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/lox/parfait/stacks"
	"gopkg.in/alecthomas/kingpin.v2"
)

func ConfigureListStacks(app *kingpin.Application, sess client.ConfigProvider) {
	var showAll bool

	cmd := app.Command("list-stacks", "List all cloudformation stacks")
	cmd.Alias("list")
	cmd.Alias("ls")

	cmd.Flag("all", "Show deleted stacks as well").
		Short('a').
		BoolVar(&showAll)

	cmd.Action(func(c *kingpin.ParseContext) error {
		var s []*cloudformation.Stack
		var err error

		if showAll {
			s, err = stacks.FindAll(cloudformation.New(sess))
		} else {
			s, err = stacks.FindAllActive(cloudformation.New(sess))
		}
		if err != nil {
			return err
		}

		fmt.Printf("%-60s %-40s %-20s\n", "NAME", "STATUS", "LAST UPDATED")
		for _, stack := range s {
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
