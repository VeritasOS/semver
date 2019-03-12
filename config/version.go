package config

var VERSION = "unknown"

const (
	ToolName                   = "semver"
	RepoBase                   = ""
	VersionFilePrefix          = ".version-"
	VersionStable              = "latest"
	HoursUntilNextUpgradeCheck = 24 * 7
	UpgradeCommandName         = "semver upgrade"
)
