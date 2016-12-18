package stacks

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/fatih/color"
)

func FormatStackStatus(s string) string {
	switch {
	case strings.HasSuffix(s, "COMPLETE") && !strings.HasPrefix(s, "DELETE"):
		return color.GreenString(s)
	case strings.Contains(s, "FAILED") || strings.Contains(s, "ROLLBACK"):
		return color.RedString(s)
	case strings.HasSuffix(s, "IN_PROGRESS"):
		return color.YellowString(s)
	}
	return s
}

func FormatStackEvent(event *cloudformation.StackEvent) string {
	descr := ""
	if event.ResourceStatusReason != nil {
		descr = fmt.Sprintf("=> %q", *event.ResourceStatusReason)
	}
	return fmt.Sprintf("%s -> %s [%s] %s",
		FormatStackStatus(*event.ResourceStatus),
		*event.LogicalResourceId,
		*event.ResourceType,
		descr,
	)
}
