package main

import (
	"code.cloudfoundry.org/cli/cf/i18n"
	"code.cloudfoundry.org/cli/cf/manifest"
	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/cf/trace"
	"code.cloudfoundry.org/cli/plugin"
	"os"
)

type localeGetter struct{}

func (l localeGetter) Locale() string {
	return "en-us"
}

func main() {
	p := new(raiPlugin)
	traceLogger := trace.NewLogger(os.Stdout, true, os.Getenv("CF_TRACE"), "")
	p.ui = terminal.NewUI(os.Stdin, os.Stdout, terminal.NewTeePrinter(os.Stdout), traceLogger)

	// Set up the manifest repository, otherwise `ReadManifest` will not work
	p.manifestRepo = manifest.NewDiskRepository()

	// Setup i18n for the creation of CF error messages created while reading manifest files (https://github.com/cloudfoundry/cli/issues/1018).
	i18n.T = i18n.Init(localeGetter{})

	p.commander = &realCommander{}

	plugin.Start(p)
}
