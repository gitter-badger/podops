package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/podops/podops"
)

// AuthCommand logs into the PodOps service and validates the token
func AuthCommand(c *cli.Context) error {
	token := c.Args().First()

	if token != "" {
		// remove old config if it exists
		if err := removeConfig(); err != nil {
			return err
		}

		// create a new client and force token verification
		cl, err := podops.NewClient(token)
		if err != nil {
			fmt.Println("\nNot authorized")
			return nil
		}
		err = cl.Store(defaultPathAndName)
		if err != nil {
			fmt.Printf("\nCould not write config. %v\n", err)
			return nil
		}

		fmt.Println("\nAuthentication successful")
	} else {
		fmt.Println("\nMissing token")
	}

	return nil
}

// LogoutCommand clears all session information
func LogoutCommand(c *cli.Context) error {
	if err := removeConfig(); err != nil {
		return err
	}

	fmt.Println("\nLogout successful")
	return nil
}
