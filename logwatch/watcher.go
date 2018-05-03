package logwatch

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	cwl "github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

type Event cwl.FilteredLogEvent

type awsApi interface {
	DescribeLogStreamsPages(input *cwl.DescribeLogStreamsInput, fn func(p *cwl.DescribeLogStreamsOutput, lastPage bool) (shouldContinue bool)) error
	FilterLogEventsPages(input *cwl.FilterLogEventsInput, fn func(p *cwl.FilterLogEventsOutput, lastPage bool) (shouldContinue bool)) error
}

type LogWatcher struct {
	LogGroup  string
	LogPrefix string
	awsApi    awsApi
}

func NewLogWatcher(awsApi awsApi, group, prefix string) *LogWatcher {
	return &LogWatcher{
		LogGroup:  group,
		LogPrefix: prefix,
		awsApi:    awsApi,
	}
}

// waitForStreams polls for steams to appear. Because of eventual consistency, this can sometimes take a while
func (lw *LogWatcher) waitForStreams(ctx context.Context, timeout time.Duration) ([]*string, error) {
	subctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var streams []*string
	var err error

	for {
		select {
		case <-time.After(2 * time.Second):
			streams, err = lw.describeStreams()
			if err != nil {
				return nil, err
			}
			if len(streams) > 0 {
				return streams, nil
			}

		case <-subctx.Done():
			return nil, subctx.Err()
		}
	}
}

func (lw *LogWatcher) describeStreams() ([]*string, error) {
	params := &cwl.DescribeLogStreamsInput{
		LogGroupName: aws.String(lw.LogGroup),
		Descending:   aws.Bool(true),
	}

	if lw.LogPrefix != "" {
		params.LogStreamNamePrefix = aws.String(lw.LogPrefix)
	}

	streams := []*string{}
	err := lw.awsApi.DescribeLogStreamsPages(params, func(page *cwl.DescribeLogStreamsOutput, lastPage bool) bool {
		for _, stream := range page.LogStreams {
			streams = append(streams, stream.LogStreamName)
		}
		return lastPage
	})

	return streams, err
}

func (lw *LogWatcher) readEventsAfter(streams []*string, ts int64, events chan *Event) (int64, error) {
	filterInput := &cwl.FilterLogEventsInput{
		LogGroupName:   aws.String(lw.LogGroup),
		LogStreamNames: streams,
		StartTime:      aws.Int64(ts + 1),
	}
	err := lw.awsApi.FilterLogEventsPages(filterInput,
		func(p *cwl.FilterLogEventsOutput, lastPage bool) (shouldContinue bool) {
			for _, event := range p.Events {
				ts = *event.Timestamp
				events <- (*Event)(event)
			}
			return lastPage
		})

	return ts, err
}

func (lw *LogWatcher) PrintEvent(ev Event) {
	fmt.Printf("%s %s\n",
		parseEventTime(*ev.Timestamp).Local().Format("2006/01/02 15:04:05"),
		*ev.Message,
	)
}

func (lw *LogWatcher) Watch(ctx context.Context, events chan *Event) error {
	streams, err := lw.waitForStreams(ctx, time.Second*30)
	if err != nil {
		return err
	}

	var after int64
	if after, err = lw.readEventsAfter(streams, after, events); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-time.After(1 * time.Second):
			if after, err = lw.readEventsAfter(streams, after, events); err != nil {
				return err
			}
		}
	}
}

func parseEventTime(millis int64) time.Time {
	return time.Unix(0, millis*int64(time.Millisecond))
}
