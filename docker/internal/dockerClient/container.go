package dockerClient

import (
	"context"
	"dock/internal/colorfmt"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func GetContainers(all bool) ([]container.Summary, error) {
	cli, err := client.NewClientWithOpts(client.WithVersion(DockerAPI), client.FromEnv)
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
		return nil, err
	}

	defer cli.Close()

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: all})
	if err != nil {
		log.Fatalf("Failed to list containers: %v", err)
		return nil, err
	}

	return containers, nil
}

func StartContainer(containerID string) error {
	cli, err := client.NewClientWithOpts(client.WithVersion(DockerAPI), client.FromEnv)
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
		return err
	}

	defer cli.Close()

	err = cli.ContainerStart(context.Background(), containerID, container.StartOptions{})
	if err != nil {
		log.Fatalf("Failed to start container: %v", err)
		return err
	}

	return nil
}

func StopContainer(containerID string) error {
	cli, err := client.NewClientWithOpts(client.WithVersion(DockerAPI), client.FromEnv)
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
		return err
	}

	defer cli.Close()

	err = cli.ContainerStop(context.Background(), containerID, container.StopOptions{})
	if err != nil {
		log.Fatalf("Failed to stop container: %v", err)
		return err
	}

	return nil
}

func ToggleContainer(containerID string) error {
	running, err := IsContainerRunning(containerID)
	if err != nil {
		log.Fatalf("Failed to check if container is running: %v", err)
		return err
	}

	if running {
		fmt.Printf("🕛 Stopping container %s...\n", containerID[:12])
		return StopContainer(containerID)
	} else {
		fmt.Printf("🕛 Starting container %s...\n", containerID[:12])
		return StartContainer(containerID)
	}
}

func IsContainerRunning(containerID string) (bool, error) {
	cli, err := client.NewClientWithOpts(client.WithVersion(DockerAPI), client.FromEnv)
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
		return false, err
	}

	defer cli.Close()

	container, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		log.Fatalf("Failed to inspect container: %v", err)
		return false, err
	}

	return container.State.Running, nil
}

func WaitForContainerToStart(containerID string) error {
	for i := 1; i <= 15; i++ {
		running, err := IsContainerRunning(containerID)
		if err != nil {
			log.Fatalf("Failed to check if container is running: %v", err)
			return err
		}

		if running {
			return nil
		}
		time.Sleep(StartupInterval)
	}

	return errors.New("🚫 Container taking too long to start. Please do it manually and try again.")
}

func GetStatusColor(status string) string {
	switch status {
	case "running":
		return colorfmt.Green
	case "exited":
		return colorfmt.Red
	case "paused":
		return colorfmt.Red
	case "restarting":
		return colorfmt.Magenta
	default:
		return colorfmt.Blue
	}
}

func GetStatusEmoji(status string) string {
	switch status {
	case "running":
		return "🟢"
	case "exited":
		return "🔴"
	case "paused":
		return "⏸"
	case "restarting":
		return "🔄"
	default:
		return "⏳"
	}
}
