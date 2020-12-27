package cli

// https://github.com/urfave/cli/blob/master/docs/v2/manual.md

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"

	"github.com/podops/podops/internal/resources"
)

// remove the local file with login credentials and other state information
func close() error {
	// remove the .po file if it exists
	f, _ := os.Stat(presetsNameAndPath)
	if f != nil {
		err := os.Remove(presetsNameAndPath)
		if err != nil {
			return err
		}
	}
	return nil
}

// PrintError formats a CLI error and prints it
func PrintError(c *cli.Context, err error) {
	msg := fmt.Sprintf("%s: %v", c.Command.Name, strings.ToLower(err.Error()))
	fmt.Println(msg)
}

// NoOpCommand is just a placeholder
func NoOpCommand(c *cli.Context) error {
	return cli.Exit(fmt.Sprintf("Command '%s' is not implemented", c.Command.Name), 0)
}

func loadResource(path string) (interface{}, string, string, error) {
	// FIXME: only local yaml is supported at the moment !

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, "", "", fmt.Errorf("Can not read file '%s'. %w", path, err)
	}

	r, kind, guid, err := resources.LoadResource(data)
	if err != nil {
		return nil, "", "", err
	}

	return r, kind, guid, nil
}

func dump(path string, doc interface{}) error {
	data, err := yaml.Marshal(doc)
	if err != nil {
		return err
	}

	ioutil.WriteFile(path, data, 0644)
	fmt.Printf("--- %s:\n\n%s\n\n", path, string(data))

	return nil
}
