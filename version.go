package main

import (
	"code.cloudfoundry.org/cli/plugin"
	"fmt"
	"github.com/blang/semver"
)

var artifactVersion string
var defaultVersion = plugin.VersionType{
	Major: 0,
	Minor: 1,
	Build: 0,
}

func pluginVersion(version string) plugin.VersionType {
	if version != "" {
		return parseVersion(version)
	}
	return defaultVersion
}

func parseVersion(version string) plugin.VersionType {
	parsedVersion, err := semver.Parse(version)
	if err != nil {
		fmt.Printf("Wrong artifactVersion specified: '%s'; using default %s\n", version, versionTypeToString(defaultVersion))
		return defaultVersion
	}
	return plugin.VersionType{Major: int(parsedVersion.Major), Minor: int(parsedVersion.Minor), Build: int(parsedVersion.Patch)}
}

func versionTypeToString(versionType plugin.VersionType) string {
	return fmt.Sprintf("%d.%d.%d", versionType.Major, versionType.Minor, versionType.Build)
}
