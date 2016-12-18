package cmd

import (
	"io/ioutil"
	"time"

	"github.com/lox/parfait/api"
	"gopkg.in/alecthomas/kingpin.v2"
)

func ConfigureUpdateStack(app *kingpin.Application, svc api.Services) {
	var stackName, tpl string
	var params []string

	cmd := app.Command("update-stack", "Update a cloudformation stack")
	cmd.Alias("update")

	cmd.Flag("file", "The cloudformation template file path").
		Short('f').
		StringVar(&tpl)

	cmd.Arg("name", "The name of the cloudformation stack").
		Required().
		StringVar(&stackName)

	cmd.Arg("params", "Parameters to the stack in Key=Val form").
		StringsVar(&params)

	cmd.Action(func(c *kingpin.ParseContext) error {
		params, err := parseStackParams(params)
		if err != nil {
			return err
		}

		b, err := ioutil.ReadFile(tpl)
		if err != nil {
			return err
		}

		ctx := api.CreateStackContext{
			Params: params,
		}

		t := time.Now()
		if err = api.CreateStack(svc.Cloudformation, stackName, string(b), ctx); err != nil {
			return err
		}

		return watchStack(svc, stackName, t)
	})
}
