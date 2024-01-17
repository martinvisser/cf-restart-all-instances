package main

import (
	"code.cloudfoundry.org/cli/cf/flags"
	"code.cloudfoundry.org/cli/plugin"
	"errors"
	"fmt"
	"strconv"
	"time"
)

var waitingAnimation = []string{
	"🕛", "🕐", "🕑", "🕒", "🕓", "🕔", "🕕", "🕖", "🕗", "🕘", "🕙", "🕚",
}

func (p *raiPlugin) restartAllInstances(cli plugin.CliConnection, args []string) error {
	applicationName := args[1]

	fc := flags.New()
	fc.NewIntFlagWithDefault("delay", "d", "Delay in seconds between restarts", 60)
	err := fc.Parse(args...)
	if err != nil {
		printHelp(cli, args)
		return err
	}

	delay := time.Duration(fc.Int("delay")) * time.Second

	app, err := cli.GetApp(applicationName)
	if err != nil {
		return err
	}

	if app.State == "stopped" {
		p.ui.Warn("Application %s is currently stopped", applicationName)
		return nil
	}

	instances := app.Instances
	p.ui.Say("Going to restart %d instance(s)", len(instances))

	for instance := range instances {
		p.ui.Say(fmt.Sprintf("Restarting instance %d", instance))

		_, err = cli.CliCommandWithoutTerminalOutput("restart-app-instance", applicationName, strconv.Itoa(instance))
		if err != nil {
			return err
		}

		err = p.postponeCheck(cli, instance, applicationName)
		if err != nil {
			return err
		}

		err = p.waitForRestart(cli, instance, applicationName, delay)
		if err != nil {
			return err
		}
		p.ui.Say("")
	}

	p.ui.Say("All instances of %s have been restarted", applicationName)

	return nil
}

func (p *raiPlugin) postponeCheck(cli plugin.CliConnection, instance int, applicationName string) error {
	for start := time.Now(); time.Since(start) < time.Duration(len(waitingAnimation))*time.Second; {
		app, err := cli.GetApp(applicationName)
		if err != nil {
			return err
		}

		state := app.Instances[instance].State
		if state == "down" || state == "starting" {
			break
		}

		for _, w := range waitingAnimation {
			fmt.Printf("\r%s Waiting until instance %d has been terminated", w, instance)
			time.Sleep(50 * time.Millisecond)
		}
	}
	return nil
}

func (p *raiPlugin) waitForRestart(cli plugin.CliConnection, instance int, applicationName string, delay time.Duration) error {
	deadline := time.Now().Add(1 * time.Minute)

	for {
		app, err := cli.GetApp(applicationName)
		if err != nil {
			return err
		}

		if time.Now().After(deadline) {
			p.ui.Failed("\nTimeout when restarting instance %d of '%s'", instance, applicationName)
			return errors.New("failed restart instance")
		} else if app.Instances[instance].State == "running" {
			p.ui.Say("\nInstance %d is restarted", instance)
			err = p.pauseForNextRestart(cli, instance, applicationName, delay)
			if err != nil {
				return err
			}
			break
		} else {
			p.animateRestart(instance)
		}
	}
	return nil
}

func (p *raiPlugin) animateRestart(instance int) {
	for start := time.Now(); time.Since(start) < time.Duration(len(waitingAnimation))*time.Second; {
		for _, w := range waitingAnimation {
			fmt.Printf("\r%s Waiting until instance %d has been restarted...", w, instance)
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func (p *raiPlugin) pauseForNextRestart(cli plugin.CliConnection, instance int, applicationName string, delay time.Duration) error {
	app, err := cli.GetApp(applicationName)
	if err != nil {
		return err
	}

	if instance < app.RunningInstances-1 {
		p.ui.Say("")
		for start := time.Now(); time.Since(start) < delay; {
			for _, w := range waitingAnimation {
				fmt.Printf("\r%s Pausing %d seconds before next instance will be restarted...", w, int64(delay.Seconds()))
				time.Sleep(250 * time.Millisecond)
			}
		}
	}
	return nil
}
