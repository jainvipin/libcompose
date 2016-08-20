package app

import (
	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/deploy"
)

func PrePlugin (v interface{}, context *cli.Context) error {
	p, ok := v.(*project.Project)
	if !ok {
		logrus.Fatalf("invalid project %v", v)
	}
	if err := deploy.PreHooks(p, context.Command.Name); err != nil {
		logrus.Fatalf("Prehook failure - Error %v", err)
	}

  return nil
}

func PostPlugin(v interface{}, context *cli.Context) error {
	p, ok := v.(*project.Project)
	if !ok {
		logrus.Fatalf("invalid project %v", v)
	}

	if err := deploy.PostHooks(p, context.Command.Name); err != nil {
		logrus.Fatalf("Posthook failure - Error %v", err)
	}
	return nil
}
