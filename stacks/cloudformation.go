package stacks

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

type cfnInterface interface {
	DescribeStacksPages(*cloudformation.DescribeStacksInput, func(*cloudformation.DescribeStacksOutput, bool) bool) error
	DescribeStackEventsPages(*cloudformation.DescribeStackEventsInput, func(*cloudformation.DescribeStackEventsOutput, bool) bool) error
	DescribeStacks(*cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error)
	CreateStack(*cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error)
	DeleteStack(*cloudformation.DeleteStackInput) (*cloudformation.DeleteStackOutput, error)
	UpdateStack(*cloudformation.UpdateStackInput) (*cloudformation.UpdateStackOutput, error)
	GetTemplate(input *cloudformation.GetTemplateInput) (*cloudformation.GetTemplateOutput, error)
	ValidateTemplate(input *cloudformation.ValidateTemplateInput) (*cloudformation.ValidateTemplateOutput, error)
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
			aws.String("CAPABILITY_NAMED_IAM"),
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

type UpdateStackContext struct {
	Params map[string]string
	Body   string
}

func Update(svc cfnInterface, name string, ctx UpdateStackContext) error {
	paramsSlice := []*cloudformation.Parameter{}
	for k, v := range ctx.Params {
		paramsSlice = append(paramsSlice, &cloudformation.Parameter{
			ParameterKey:   aws.String(k),
			ParameterValue: aws.String(v),
		})
	}

	if ctx.Body == "" {
		log.Printf("Reading previous template")
		resp, err := svc.GetTemplate(&cloudformation.GetTemplateInput{
			StackName: aws.String(name),
		})
		if err != nil {
			return err
		}
		ctx.Body = *resp.TemplateBody
	}

	// validate to get parameters
	validate, err := svc.ValidateTemplate(&cloudformation.ValidateTemplateInput{
		TemplateBody: aws.String(ctx.Body),
	})
	if err != nil {
		return err
	}

	// use previous values for any missing params
	for _, param := range validate.Parameters {
		if _, hasParam := ctx.Params[*param.ParameterKey]; !hasParam {
			log.Printf("Using previous value for %s", *param.ParameterKey)
			paramsSlice = append(paramsSlice, &cloudformation.Parameter{
				ParameterKey:     param.ParameterKey,
				UsePreviousValue: aws.Bool(true),
			})
		}
	}

	_, err = svc.UpdateStack(&cloudformation.UpdateStackInput{
		StackName: aws.String(name),
		Capabilities: []*string{
			aws.String("CAPABILITY_IAM"),
			aws.String("CAPABILITY_NAMED_IAM"),
		},
		Parameters:   paramsSlice,
		TemplateBody: aws.String(ctx.Body),
	})
	return err
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

func IsFailed(svc cfnInterface, name string) (bool, error) {
	resp, err := svc.DescribeStacks(&cloudformation.DescribeStacksInput{
		StackName: aws.String(name),
	})
	if err != nil {
		return true, err
	}

	if len(resp.Stacks) != 1 {
		return true, fmt.Errorf("Expected 1 stack, got %d", len(resp.Stacks))
	}

	log.Printf("%#v", *resp.Stacks[0].StackStatus)

	outputs := map[string]string{}
	for _, output := range resp.Stacks[0].Outputs {
		outputs[*output.OutputKey] = *output.OutputValue
	}

	return false, nil
}