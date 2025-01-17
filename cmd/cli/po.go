package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	a "github.com/podops/podops/apiv1"
	cmd "github.com/podops/podops/cmd/cli/commands"
	"github.com/urfave/cli/v2"
)

const (
	cmdLineName = "po"
)

func main() {

	// initialize CLI
	app := &cli.App{
		Name:    cmdLineName,
		Version: a.VersionString,
		Usage:   fmt.Sprintf("PodOps: Podcast Operations CLI (%s)", a.Version),
		Action: func(c *cli.Context) error {
			fmt.Println(globalHelpText)
			return nil
		},
		Commands: setupCommands(),
	}

	sort.Sort(cli.FlagsByName(app.Flags))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func setupCommands() []*cli.Command {
	c := []*cli.Command{
		// Basic Commands
		{
			Name:     "list",
			Usage:    "List all shows/productions",
			Category: cmd.BasicCmdGroup,
			Action:   cmd.ListProductionsCommand,
		},
		{
			Name:      "set",
			Usage:     "List the current show/production, switch to another show/production",
			UsageText: setUsageText,
			Category:  cmd.BasicCmdGroup,
			Action:    cmd.SetProductionCommand,
		},
		{
			Name:      "new-show",
			Usage:     "Setup a new show/production",
			UsageText: "new-show NAME",
			Category:  cmd.BasicCmdGroup,
			Action:    cmd.NewProductionCommand,
			Flags:     newShowFlags(),
		},
		{
			Name:      "template",
			Usage:     "Create a resource template with default values",
			UsageText: "template [show|episode] NAME",
			Category:  cmd.BasicCmdGroup,
			Action:    cmd.TemplateCommand,
			Flags:     templateFlags(),
		},

		// Settings
		{
			Name:      "auth",
			Usage:     "Login to the service",
			UsageText: "auth TOKEN",
			Category:  cmd.SettingsCmdGroup,
			Action:    cmd.AuthCommand,
		},
		{
			Name:     "logout",
			Usage:    "Logout and clear all session information",
			Category: cmd.SettingsCmdGroup,
			Action:   cmd.LogoutCommand,
		},
		// Show/production commands
		{
			Name:      "get",
			Usage:     "Lists a single resource/a collection of resources",
			UsageText: getUsageText,
			Category:  cmd.ShowCmdGroup,
			Action:    cmd.ListResourcesCommand,
			//Flags:     createFlags(),
		},
		{
			Name:      "create",
			Usage:     "Create a resource from a file, directory or URL",
			UsageText: "create FILENAME",
			Category:  cmd.ShowCmdGroup,
			Action:    cmd.CreateCommand,
			Flags:     createFlags(),
		},
		{
			Name:      "update",
			Usage:     "Update a resource from a file, directory or URL",
			UsageText: "update FILENAME",
			Category:  cmd.ShowCmdGroup,
			Action:    cmd.UpdateCommand,
			Flags:     createFlags(),
		},
		{
			Name:      "upload",
			Usage:     "Upload an asset from a file",
			UsageText: "upload FILENAME",
			Category:  cmd.ShowCmdGroup,
			Action:    cmd.UploadCommand,
			Flags:     createFlags(),
		},
		{
			Name:      "build",
			Usage:     "Start a new build",
			UsageText: "po build",
			Category:  cmd.ShowMgmtCmdGroup,
			Action:    cmd.BuildCommand,
		},
		{
			Name:      "delete",
			Usage:     "Delete a resource",
			UsageText: "po delete [show|episode] NAME",
			Category:  cmd.ShowMgmtCmdGroup,
			Action:    cmd.DeleteResourcesCommand,
		},
	}
	return c
}

func globalFlags() []cli.Flag {
	f := []cli.Flag{
		&cli.BoolFlag{
			Name:    "debug",
			Usage:   "Prints debug information",
			Aliases: []string{"d"},
		},
	}
	return f
}

func newShowFlags() []cli.Flag {
	f := []cli.Flag{
		&cli.StringFlag{
			Name:    "title",
			Usage:   "Show title",
			Aliases: []string{"t"},
		},
		&cli.StringFlag{
			Name:    "summary",
			Usage:   "Show summary",
			Aliases: []string{"s"},
		},
	}
	return f
}

func createFlags() []cli.Flag {
	f := []cli.Flag{
		&cli.BoolFlag{
			Name:    "force",
			Usage:   "Force create/update/upload",
			Aliases: []string{"f"},
		},
	}
	return f
}

func templateFlags() []cli.Flag {
	f := []cli.Flag{
		/*
			&cli.StringFlag{
				Name:    "name",
				Usage:   "Resource name",
				Aliases: []string{"n"},
			},
		*/
		&cli.StringFlag{
			Name:    "parent",
			Usage:   "Parent resource name",
			Aliases: []string{"p"},
		},
		&cli.StringFlag{
			Name:    "guid",
			Usage:   "Resource GUID",
			Aliases: []string{"id"},
		},
		&cli.StringFlag{
			Name:    "parentid",
			Usage:   "Parent resource GUID",
			Aliases: []string{"pid"},
		},
	}
	return f
}

//
// all the help texts used in the CLI
//
const (
	globalHelpText = `PodOps: Podcast Operations Client

This client tool helps you to create and produce podcasts.
It also includes administrative commands for managing your live podcasts.

To see the full list of supported commands, run 'po help'`

	setUsageText = `set [NAME]

	 # Display the current show/production
	 po set
	 
	 # Set the current show/production
	 po set [NAME]`

	getUsageText = `get [RESOURCE]

	 # List all resources
	 po get

	 # List all resources of a type
	 po get [show|episode]

	 # Show details about a resource
	 po get [show|episode] NAME`
)
