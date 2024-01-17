package main

import (
	"code.cloudfoundry.org/cli/plugin"
)

var commands = []plugin.Command{
	commandRestartAllInstance,
}

var commandRestartAllInstance = plugin.Command{
	Name:     "restart-all-instances",
	HelpText: "Restart all instances with a configurable delay",
	UsageDetails: plugin.Usage{
		Usage: "cf restart-all-instances <application-name> [-d delay-in-seconds]",
		Options: map[string]string{
			"d, --delay": "Delay in seconds",
		},
	},
}

type rocsCommand struct {
	f             func(*raiPlugin, plugin.CliConnection, []string) error
	argCount      int
	loginRequired bool
}

var functions = map[string]rocsCommand{
	commandRestartAllInstance.Name: {f: (*raiPlugin).restartAllInstances, argCount: 2, loginRequired: true},
	"CLI-MESSAGE-UNINSTALL":        {f: (*raiPlugin).uninstallHook, argCount: 1, loginRequired: false},
}
