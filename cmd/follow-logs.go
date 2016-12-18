package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	logs "github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/lox/parfait/api"
	"gopkg.in/alecthomas/kingpin.v2"
)

func ConfigureFollowLogs(app *kingpin.Application, svc api.Services) {
	var logGroup, prefix string

	cmd := app.Command("follow-logs", "Follow a cloudwatch log group")
	cmd.Alias("logs")

	cmd.Flag("log-group", "The cloudwatch logs group to follow").
		Short('g').
		StringVar(&logGroup)

	cmd.Flag("prefix", "Filter log streams by this prefix").
		Short('p').
		StringVar(&prefix)

	cmd.Action(func(c *kingpin.ParseContext) error {
		params := &logs.DescribeLogStreamsInput{
			LogGroupName:        aws.String(logGroup),
			LogStreamNamePrefix: aws.String(prefix),
			Descending:          aws.Bool(true),
		}

		streams := []*string{}
		err := svc.Logs.DescribeLogStreamsPages(params, func(page *logs.DescribeLogStreamsOutput, lastPage bool) bool {
			for _, stream := range page.LogStreams {
				streams = append(streams, stream.LogStreamName)
			}
			return lastPage
		})

		if err != nil {
			return err
		}

		filterInput := &logs.FilterLogEventsInput{
			LogGroupName:   aws.String(logGroup),
			LogStreamNames: streams,
		}

		err = svc.Logs.FilterLogEventsPages(filterInput, func(p *logs.FilterLogEventsOutput, lastPage bool) (shouldContinue bool) {
			for _, event := range p.Events {
				printLogEvent(event)
			}
			return lastPage
		})

		if err != nil {
			return err
		}

		return nil
	})
}

func printLogEvent(ev *logs.FilteredLogEvent) {
	name := *ev.LogStreamName
	if len(name) > 40 {
		name = name[:37] + "..."
	}

	fmt.Printf("%-20d %-42s %s\n", *ev.Timestamp, name, *ev.Message)
}
