package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/lox/parfait/api"
	"gopkg.in/alecthomas/kingpin.v2"
)

func ConfigureCreateStack(app *kingpin.Application, svc api.Services) {
	var stackName, tpl string
	var params []string

	cmd := app.Command("create-stack", "Create a cloudformation stack")
	cmd.Alias("create")
	cmd.Flag("stack-name", "The name of the cloudformation stack").
		Short('n').
		StringVar(&stackName)

	cmd.Flag("file", "The cloudformation template").
		Short('f').
		StringVar(&tpl)

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
		log.Printf("%#v", ctx)

		if err = api.CreateStack(svc.Cloudformation, stackName, string(b), ctx); err != nil {
			return err
		}

		return watchStack(svc, stackName)
	})
}

func parseStackParams(rawParams []string) (map[string]string, error) {
	params := map[string]string{}
	for _, arg := range rawParams {
		parts := strings.Split(arg, "=")
		if len(parts) != 2 {
			return params, fmt.Errorf("Failed to parse parameter %q", arg)
		}
		params[parts[0]] = parts[1]
	}
	return params, nil
}
