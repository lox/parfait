package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/lox/parfait/api"
	"gopkg.in/alecthomas/kingpin.v2"
)

func ConfigureCreateStack(app *kingpin.Application, svc api.Services) {
	var stackName, tpl string
	var params []string
	var tplURL *url.URL

	cmd := app.Command("create-stack", "Create a cloudformation stack")
	cmd.Alias("create")

	cmd.Flag("file", "The file path to a cloudformation template").
		Short('f').
		StringVar(&tpl)

	cmd.Flag("url", "The url to a cloudformation template").
		Short('u').
		URLVar(&tplURL)

	cmd.Arg("stack-name", "The name of the cloudformation stack").
		StringVar(&stackName)

	cmd.Arg("params", "Parameters to the stack in Key=Val form").
		StringsVar(&params)

	cmd.Action(func(c *kingpin.ParseContext) error {

		// validate params
		if tplURL == nil && tpl == "" {
			return errors.New("Must provide either --url or --file")
		} else if tplURL != nil && tpl != "" {
			return errors.New("Can't provide both --url and --file")
		}

		params, err := parseStackParams(params)
		if err != nil {
			return err
		}

		var b []byte

		if tplURL != nil {
			if b, err = readURL(tplURL); err != nil {
				return err
			}
		} else {
			if b, err = ioutil.ReadFile(tpl); err != nil {
				return err
			}
		}

		ctx := api.CreateStackContext{
			Params: params,
		}

		if err = api.CreateStack(svc.Cloudformation, stackName, string(b), ctx); err != nil {
			return err
		}

		return watchStack(svc, stackName)
	})
}

func readURL(u *url.URL) ([]byte, error) {
	response, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Response status was %s", response.Status)
	}

	return ioutil.ReadAll(response.Body)
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
