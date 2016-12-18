package logwatch

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

type Event cloudwatchlogs.FilteredLogEvent

type cloudwatchApi interface {
	DescribeLogStreamsPages(input *cloudwatchlogs.DescribeLogStreamsInput, fn func(p *cloudwatchlogs.DescribeLogStreamsOutput, lastPage bool) (shouldContinue bool)) error
	FilterLogEventsPages(input *cloudwatchlogs.FilterLogEventsInput, fn func(p *cloudwatchlogs.FilterLogEventsOutput, lastPage bool) (shouldContinue bool)) error
}

type LogWatcher struct {
	LogGroup  string
	LogPrefix string

	api cloudwatchApi
}

func NewLogWatcher(group, prefix string, api cloudwatchApi) *LogWatcher {
	return &LogWatcher{
		LogGroup:  group,
		LogPrefix: prefix,
		api:       api,
	}
}

func (lw *LogWatcher) pollStreams(ctx context.Context) ([]*string, error) {
	var streams []*string
	var err error

	for {
		select {
		case <-time.After(1 * time.Second):
			streams, err = lw.describeStreams()
			if err != nil {
				return nil, err
			}
			if len(streams) > 0 {
				return streams, nil
			}

		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (lw *LogWatcher) describeStreams() ([]*string, error) {
	params := &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName:        aws.String(lw.LogGroup),
		LogStreamNamePrefix: aws.String(lw.LogPrefix),
		Descending:          aws.Bool(true),
	}

	streams := []*string{}
	err := lw.api.DescribeLogStreamsPages(params, func(page *cloudwatchlogs.DescribeLogStreamsOutput, lastPage bool) bool {
		for _, stream := range page.LogStreams {
			streams = append(streams, stream.LogStreamName)
		}
		return lastPage
	})

	return streams, err
}

func (lw *LogWatcher) readEventsAfter(streams []*string, ts int64, events chan *Event) (int64, error) {
	filterInput := &cloudwatchlogs.FilterLogEventsInput{
		LogGroupName:   aws.String(lw.LogGroup),
		LogStreamNames: streams,
		StartTime:      aws.Int64(ts + 1),
	}
	err := lw.api.FilterLogEventsPages(filterInput,
		func(p *cloudwatchlogs.FilterLogEventsOutput, lastPage bool) (shouldContinue bool) {
			for _, event := range p.Events {
				ts = *event.Timestamp
				events <- (*Event)(event)
			}
			return lastPage
		})

	return ts, err
}

func (lw *LogWatcher) PrintEvent(ev Event) {
	name := strings.TrimPrefix(*ev.LogStreamName, lw.LogPrefix)
	if len(name) > 50 {
		name = name[:47] + "..."
	}

	fmt.Printf("%-20s %-52s %s\n",
		parseEventTime(*ev.Timestamp).Format(time.Stamp),
		name,
		*ev.Message,
	)
}

func (lw *LogWatcher) Watch(ctx context.Context, events chan *Event) error {
	subctx, _ := context.WithTimeout(ctx, time.Second*5)
	streams, err := lw.pollStreams(subctx)
	if err != nil {
		return err
	}

	for _, stream := range streams {
		log.Printf("Found stream %s", *stream)
	}

	var after int64
	if after, err = lw.readEventsAfter(streams, after, events); err != nil {
		return err
	}

	for {
		select {
		case <-time.After(1 * time.Second):
			if after, err = lw.readEventsAfter(streams, after, events); err != nil {
				return err
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

func parseEventTime(millis int64) time.Time {
	return time.Unix(0, millis*int64(time.Millisecond))
}
