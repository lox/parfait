package logwatch

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

type cloudwatchStub struct {
	LogStreams []*cloudwatchlogs.LogStream
	Events     []*cloudwatchlogs.FilteredLogEvent
}

func (cw *cloudwatchStub) DescribeLogStreamsPages(input *cloudwatchlogs.DescribeLogStreamsInput, fn func(p *cloudwatchlogs.DescribeLogStreamsOutput, lastPage bool) (shouldContinue bool)) error {
	fn(&cloudwatchlogs.DescribeLogStreamsOutput{
		LogStreams: cw.LogStreams,
	}, true)
	return nil
}

func (cw *cloudwatchStub) FilterLogEventsPages(input *cloudwatchlogs.FilterLogEventsInput, fn func(p *cloudwatchlogs.FilterLogEventsOutput, lastPage bool) (shouldContinue bool)) error {
	fn(&cloudwatchlogs.FilterLogEventsOutput{
		Events: cw.Events,
	}, true)
	return nil
}

func TestWatchingSimpleLog(t *testing.T) {
	cw := &cloudwatchStub{
		LogStreams: []*cloudwatchlogs.LogStream{
			&cloudwatchlogs.LogStream{
				Arn:           aws.String("arn::llamas"),
				LogStreamName: aws.String("llamas"),
			},
		},
		Events: []*cloudwatchlogs.FilteredLogEvent{
			&cloudwatchlogs.FilteredLogEvent{
				Message:   aws.String("my event 1"),
				Timestamp: aws.Int64(1),
			},
			&cloudwatchlogs.FilteredLogEvent{
				Message:   aws.String("my event 2"),
				Timestamp: aws.Int64(2),
			},
			&cloudwatchlogs.FilteredLogEvent{
				Message:   aws.String("my event 3"),
				Timestamp: aws.Int64(3),
			},
		},
	}

	w := &LogWatcher{
		LogGroup:  "myGroup",
		LogPrefix: "myPrefix",
		awsApi:    cw,
	}
	events := make(chan *Event)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		err := w.Watch(ctx, events)
		if err != context.Canceled {
			t.Fatal(err)
		}
		close(events)
	}()

	for i := 0; i < 3; i++ {
		ev := <-events
		if ev == nil {
			t.Fatalf("Expected non-nil event for event %d", i)
		}
	}

	cancel()
	wg.Wait()
}
