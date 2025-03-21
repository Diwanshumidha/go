package dockerClient

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/docker/docker/client"
)

func IsDockerRunning() bool {
	cli, err := client.NewClientWithOpts(client.WithVersion(DockerAPI), client.FromEnv)
	if err != nil {
		return false
	}
	defer cli.Close()

	_, err = cli.Ping(context.Background())
	return err == nil
}

func StartDocker() error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux", "darwin":
		cmd = exec.Command("sudo", "systemctl", "start", "docker")
	case "windows":
		cmd = exec.Command("powershell", "start_docker")
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start Docker: %v", err)
	}


	return nil
}

func MakeSureDockerIsRunning() error {
	if IsDockerRunning() {
		fmt.Println("✅ Docker is already running")
		return nil
	}

	fmt.Println("⚠️  Docker is not running. Attempting to start...")

	err := StartDocker()
	if err != nil {
		return fmt.Errorf("❌ Attempt Failed to start Docker: %v\n", err)
	} else {
		fmt.Println("⏳ Waiting for Docker to start...")
	}


	for i := 1; i <= MaxStartupAttempts; i++ {
		time.Sleep(StartupInterval)
		if IsDockerRunning() {
			fmt.Println("✅ Docker is now running")
			return nil
		}
	}

	return errors.New("🚫 Docker taking too long to start. Please start Docker manually and try again.")
}
