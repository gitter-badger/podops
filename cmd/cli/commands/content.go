package commands

import (
	"fmt"

	"github.com/fupas/commons/pkg/util"
	a "github.com/podops/podops/apiv1"
	"github.com/urfave/cli/v2"
)

// TemplateCommand creates a resource template with all default values
func TemplateCommand(c *cli.Context) error {
	template := c.Args().First()
	if template != a.ResourceShow && template != a.ResourceEpisode {
		fmt.Println(fmt.Sprintf("\nDon't know how to create '%s'", template))
		return nil
	}

	name := "NAME"
	if c.NArg() == 2 {
		name = c.Args().Get(1)
	}
	// extract flags or set defaults
	guid := c.String("id")
	if guid == "" {
		guid, _ = util.ShortUUID()
	}
	parent := c.String("parent")
	if parent == "" {
		parent = "PARENT-NAME"
	}
	parentGUID := c.String("parentid")
	if parentGUID == "" {
		parentGUID = "PARENT-ID"
	}

	// create the yamls
	if template == "show" {

		show := a.DefaultShow(name, "TITLE", "SUMMARY", guid, a.DefaultPortalEndpoint, a.DefaultCDNEndpoint)
		err := dump(fmt.Sprintf("show-%s.yaml", guid), show)
		if err != nil {
			printError(c, err)
			return nil
		}
	} else {

		episode := a.DefaultEpisode(name, parent, guid, parentGUID, a.DefaultPortalEndpoint, a.DefaultCDNEndpoint)
		err := dump(fmt.Sprintf("episode-%s.yaml", guid), episode)
		if err != nil {
			printError(c, err)
			return nil
		}
	}

	return nil
}

// CreateCommand creates a resource from a file, directory or URL
func CreateCommand(c *cli.Context) error {

	if c.NArg() != 1 {
		return fmt.Errorf("wrong number of arguments: expected 1, got %d", c.NArg())
	}
	path := c.Args().First()
	force := c.Bool("force")

	r, kind, guid, err := loadResource(path)
	if err != nil {
		return err
	}

	_, err = client.CreateResource(kind, guid, force, r)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("created resource %s-%s", kind, guid))
	return nil
}

// UpdateCommand updates a resource from a file, directory or URL
func UpdateCommand(c *cli.Context) error {

	if c.NArg() != 1 {
		return fmt.Errorf("wrong number of arguments: expected 1, got %d", c.NArg())
	}
	path := c.Args().First()
	force := c.Bool("force")

	r, kind, guid, err := loadResource(path)
	if err != nil {
		return err
	}

	_, err = client.UpdateResource(kind, guid, force, r)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Updated resource %s-%s", kind, guid))
	return nil
}

// BuildCommand starts a new build of the feed
func BuildCommand(c *cli.Context) error {

	// FIXME support the 'NAME' option

	build, err := client.Build(client.GUID)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Build production '%s' successful.\nAccess the feed at %s", client.GUID, build.FeedAliasURL))
	return nil
}

// UploadCommand uploads an asset from a file
func UploadCommand(c *cli.Context) error {

	if c.NArg() != 1 {
		return fmt.Errorf("wrong number of arguments: expected 1, got %d", c.NArg())
	}
	name := c.Args().First()
	force := c.Bool("force")

	err := client.Upload(name, force)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Uploaded '%s'", name))
	return nil
}
