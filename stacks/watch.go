package stacks

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/lox/parfait/stacks/poller"
)

func Watch(cfn cfnInterface, stackName string, f func(event *cloudformation.StackEvent)) error {
	err := poller.UntilCreatedOrUpdated(cfn, stackName, f)
	if err != nil {
		return err
	}

	if err = GetError(cfn, stackName); err != nil {
		return err
	}

	outputs, err := Outputs(cfn, stackName)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Printf("%-20s %-80s\n", "KEY", "VALUE")
	for k, v := range outputs {
		fmt.Printf("%-20s %-80s\n", k, v)
	}

	return nil
}
