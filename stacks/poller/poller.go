package poller

import (
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
)

type cfnInterface interface {
	DescribeStackEventsPages(*cfn.DescribeStackEventsInput, func(*cfn.DescribeStackEventsOutput, bool) bool) error
}

type Poller struct {
	StackName string
	awsApi    cfnInterface
}

func NewPoller(api cfnInterface, stackName string) *Poller {
	return &Poller{
		StackName: stackName,
		awsApi:    api,
	}
}

func (p *Poller) Poll(condition EndCondition, f func(e *cfn.StackEvent)) error {
	lastSeen := time.Time{}

	for {
		events, err := p.getEvents(lastSeen)
		if err != nil {
			return err
		}

		for i := len(events) - 1; i >= 0; i-- {
			if events[i].Timestamp.After(lastSeen) {
				f(events[i])
				lastSeen = *events[i].Timestamp
			}
		}

		if len(events) > 0 {
			t, err := condition(p.StackName, events[0])
			if err != nil {
				return err
			}
			if t {
				break
			}
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}

// getEvents returns all events after a given time in reverse chronological order
func (p *Poller) getEvents(after time.Time) (events []*cfn.StackEvent, err error) {
	params := &cfn.DescribeStackEventsInput{
		StackName: aws.String(p.StackName),
	}

	err = p.awsApi.DescribeStackEventsPages(params, func(page *cfn.DescribeStackEventsOutput, last bool) bool {
		for _, event := range page.StackEvents {
			if !event.Timestamp.After(after) {
				return true
			}
			events = append(events, event)

			// stop once we hit the most recent User Initiated event
			if event.ResourceStatusReason != nil &&
				*event.ResourceStatusReason == `User Initiated` {
				return true
			}
		}
		return last
	})

	return
}

func UntilCreatedOrUpdated(api cfnInterface, stackName string, f func(e *cfn.StackEvent)) error {
	return NewPoller(api, stackName).Poll(isCreatedOrUpdated, f)
}

func UntilDeleted(api cfnInterface, stackName string, f func(e *cfn.StackEvent)) error {
	err := NewPoller(api, stackName).Poll(isDeleted, f)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// 400: ValidationError: Stack does not exist
			if awsErr.Code() == "ValidationError" {
				return nil
			}
		}
	}
	return err
}

type EndCondition func(stackName string, ev *cfn.StackEvent) (bool, error)

func isCreatedOrUpdated(stackName string, ev *cfn.StackEvent) (bool, error) {
	if *ev.LogicalResourceId == stackName {
		switch *ev.ResourceStatus {
		case cfn.ResourceStatusUpdateComplete,
			cfn.ResourceStatusCreateComplete,
			cfn.ResourceStatusUpdateFailed,
			cfn.ResourceStatusCreateFailed,
			cfn.StackStatusRollbackComplete,
			cfn.StackStatusRollbackFailed,
			cfn.StackStatusUpdateRollbackComplete,
			cfn.StackStatusUpdateRollbackFailed:
			var err error
			if ev.ResourceStatusReason != nil {
				err = errors.New(*ev.ResourceStatusReason)
			}
			return true, err
		}
	}
	return false, nil
}

func isDeleted(stackName string, ev *cfn.StackEvent) (bool, error) {
	if *ev.LogicalResourceId == stackName {
		switch *ev.ResourceStatus {
		case cfn.ResourceStatusDeleteComplete:
			return true, nil
		case cfn.ResourceStatusDeleteFailed:
			return true, errors.New(*ev.ResourceStatusReason)
		}
	}
	return false, nil
}
