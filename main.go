package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"gtithub.com/jgfranco17/lazyfile/cli/core"
)

const (
	projectName        = "lazyfile"
	projectDescription = "LazyFile: file navigation in the terminal, made easy"
)

var (
	logger  *log.Logger
	version string = "0.0.0-dev"
)

func main() {
	command := core.NewCommandRegistry(projectName, projectDescription, version, logger)
	commandsList := []*cobra.Command{
		core.CommandListFiles(),
	}
	command.RegisterCommands(commandsList)

	err := command.Execute()
	if err != nil {
		log.Error(err.Error())
	}
}
