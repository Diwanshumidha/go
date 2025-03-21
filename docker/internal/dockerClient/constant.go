package dockerClient

import "time"

const (
	DockerAPI          = "1.46"          // Docker API version
	MaxStartupAttempts = 5               // Max retries to start Docker
	StartupInterval    = 2 * time.Second // Time between retries
)
