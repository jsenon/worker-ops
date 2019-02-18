package config

// Global Config of the application
const (
	Service     = "worker-ops"
	Description = "Worker Application Reporter"
)

// Dynamic version retrieve with ldflags
// Version represent version of application
var Version string

// GitCommit represent git commit
var GitCommit string

// BuildDate represent date of build
var BuildDate string
