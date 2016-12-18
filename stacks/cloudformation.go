package stacks

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

type cfnInterface interface {
	DescribeStacksPages(*cloudformation.DescribeStacksInput, func(*cloudformation.DescribeStacksOutput, bool) bool) error
	DescribeStackEventsPages(*cloudformation.DescribeStackEventsInput, func(*cloudformation.DescribeStackEventsOutput, bool) bool) error
	DescribeStacks(*cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error)
	CreateStack(*cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error)
	DeleteStack(*cloudformation.DeleteStackInput) (*cloudformation.DeleteStackOutput, error)
}

func FindAll(svc cfnInterface) (stacks []*cloudformation.Stack, err error) {
	err = svc.DescribeStacksPages(nil, func(page *cloudformation.DescribeStacksOutput, last bool) bool {
		for _, s := range page.Stacks {
			stacks = append(stacks, s)
		}
		return last
	})
	return
}

func FindAllActive(svc cfnInterface) (stacks []*cloudformation.Stack, err error) {
	err = svc.DescribeStacksPages(nil, func(page *cloudformation.DescribeStacksOutput, last bool) bool {
		for _, s := range page.Stacks {
			if *s.StackStatus != "DELETE_COMPLETE" {
				stacks = append(stacks, s)
			}
		}
		return last
	})
	return
}

func FindByName(svc cfnInterface, stackName string) (stacks []*cloudformation.Stack, err error) {
	filter := &cloudformation.DescribeStacksInput{
		StackName: &stackName,
	}

	err = svc.DescribeStacksPages(filter, func(page *cloudformation.DescribeStacksOutput, last bool) bool {
		for _, s := range page.Stacks {
			if *s.StackStatus != "DELETE_COMPLETE" {
				stacks = append(stacks, s)
			}
		}
		return last
	})
	return
}

type CreateStackContext struct {
	Params          map[string]string
	Body            string
	DisableRollback bool
}

func Create(svc cfnInterface, name string, ctx CreateStackContext) error {
	paramsSlice := []*cloudformation.Parameter{}
	for k, v := range ctx.Params {
		paramsSlice = append(paramsSlice, &cloudformation.Parameter{
			ParameterKey:   aws.String(k),
			ParameterValue: aws.String(v),
		})
	}

	_, err := svc.CreateStack(&cloudformation.CreateStackInput{
		StackName: aws.String(name),
		Capabilities: []*string{
			aws.String("CAPABILITY_IAM"),
		},
		DisableRollback: aws.Bool(ctx.DisableRollback),
		Parameters:      paramsSlice,
		TemplateBody:    aws.String(ctx.Body),
	})
	if err != nil {
		return err
	}
	return nil
}

func Delete(svc cfnInterface, name string) error {
	_, err := svc.DeleteStack(&cloudformation.DeleteStackInput{
		StackName: &name,
	})

	return err
}

func Outputs(svc cfnInterface, name string) (map[string]string, error) {
	resp, err := svc.DescribeStacks(&cloudformation.DescribeStacksInput{
		StackName: aws.String(name),
	})

	if err != nil {
		return nil, err
	}

	if len(resp.Stacks) != 1 {
		return nil, fmt.Errorf("Expected 1 stack, got %d", len(resp.Stacks))
	}

	outputs := map[string]string{}
	for _, output := range resp.Stacks[0].Outputs {
		outputs[*output.OutputKey] = *output.OutputValue
	}

	return outputs, nil
}
