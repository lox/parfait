package stacks

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/lox/parfait/stacks/poller"
)

func Watch(cfn cfnInterface, stackName string, f func(event *cloudformation.StackEvent)) error {
	err := poller.UntilCreatedOrUpdated(cfn, stackName, f)
	if err != nil {
		return err
	}

	if failed, _ := IsFailed(cfn, stackName); failed {
		return errors.New("Stack has failed")
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
