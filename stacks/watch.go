package stacks

import (
	"log"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/lox/parfait/stacks/poller"
)

func Watch(cfn cfnInterface, stackName string, f func(event *cloudformation.StackEvent)) error {
	err := poller.UntilCreatedOrUpdated(cfn, stackName, f)
	if err != nil {
		return err
	}

	outputs, err := Outputs(cfn, stackName)
	if err != nil {
		return err
	}

	for k, v := range outputs {
		log.Printf("Stack Output: %s = %s", k, v)
	}

	return nil
}
