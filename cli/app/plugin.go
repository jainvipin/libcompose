package app

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/libcompose/deploy"
	"github.com/docker/libcompose/project"
)

func pre_plugin(p *project.Project, context *cli.Context) error {
	cliLabels := ""

	if cliLabels != "" {
		if err := deploy.PopulateEnvLabels(p, cliLabels); err != nil {
			logrus.Fatalf("Unable to insert environment labels. Error %v", err)
		}
	}

	if err := deploy.PreHooks(p, context.Command.Name); err != nil {
		logrus.Fatalf("Unable to generate network labels. Error %v", err)
	}

	return nil
}

func post_plugin(p *project.Project, context *cli.Context) error {
	if err := deploy.PostHooks(p, context.Command.Name); err != nil {
		logrus.Fatalf("Unable to populate dns entries. Error %v", err)
	}
	return nil
}
