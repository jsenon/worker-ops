package config

// Global Config of the application
const (
	Service     = "worker-ops"
	Description = "Worker Application Reporter"
)

// Dynamic version retrieve with ldflags
var Version string
var GitCommit string
var BuildDate string
