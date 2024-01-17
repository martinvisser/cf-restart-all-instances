package main

import (
	"code.cloudfoundry.org/cli/cf/manifest"
	"code.cloudfoundry.org/cli/plugin"
	"os"

	"code.cloudfoundry.org/cli/cf/terminal"
)

type raiPlugin struct {
	ui           terminal.UI
	commander    commander
	manifestRepo manifest.Repository
}

var exit = os.Exit

func (p *raiPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name:    "Restart All Instances plugin",
		Version: pluginVersion(artifactVersion),
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 51,
			Build: 0,
		},
		Commands: commands,
	}
}

func (p *raiPlugin) Run(cli plugin.CliConnection, args []string) {
	if command, ok := functions[args[0]]; ok {
		if len(args) < command.argCount {
			printHelp(cli, args)
			exit(1)
		} else if loggedIn, _ := cli.IsLoggedIn(); command.loginRequired && !loggedIn {
			p.ui.Failed("Not logged in")
			exit(1)
		} else {
			err := command.f(p, cli, args)
			if err != nil {
				p.ui.Failed(err.Error())
				exit(1)
			}
		}
	} else {
		p.ui.Say("Command '%s' not found", args[0])
		printHelp(cli, args)
		exit(1)
	}
}

func (p *raiPlugin) uninstallHook(_ plugin.CliConnection, _ []string) error {
	p.ui.Say("Sorry to see you GO! 👋")
	return nil
}

func printHelp(cli plugin.CliConnection, args []string) {
	_, _ = cli.CliCommand("help", args[0])
}
