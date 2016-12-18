package cmd

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/lox/parfait/logwatch"
	"gopkg.in/alecthomas/kingpin.v2"
)

func ConfigureFollowLogs(app *kingpin.Application, sess client.ConfigProvider) {
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
		watcher := logwatch.NewLogWatcher(cloudwatchlogs.New(sess), logGroup, prefix)
		events := make(chan *logwatch.Event)

		go func() {
			for event := range events {
				watcher.PrintEvent(*event)
			}
		}()

		return watcher.Watch(context.Background(), events)
	})
}
