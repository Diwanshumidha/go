package service

import (
	"dock/internal/dockerClient"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/docker/docker/api/types/container"
)

func ToggleContainerUsingPrompt() {
	if err := dockerClient.MakeSureDockerIsRunning(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	containers, err := dockerClient.GetContainers(true)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(containers) == 0 {
		fmt.Println("❌ No containers found.")
		os.Exit(1)
	}

	// Items := make([]string, len(containers))
	options := make([]huh.Option[string], len(containers))
	for i, container := range containers {
		emoji := dockerClient.GetStatusEmoji(container.State)
		options[i] = huh.Option[string]{Value: container.ID, Key: fmt.Sprintf("%s %s", emoji, container.Names[0][1:])}
	}

	var selectedContainerID string
	err = huh.NewSelect[string]().Title("Select a container to toggle").Options(options...).Value(&selectedContainerID).Run()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := dockerClient.ToggleContainer(selectedContainerID); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func ToggleContainerByName(name string) {
	if err := dockerClient.MakeSureDockerIsRunning(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	containers, err := dockerClient.GetContainers(true)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(containers) == 0 {
		fmt.Println("❌ No containers found.")
		os.Exit(1)
	}

	if name == "" {
		log.Fatal("Container name cannot be empty")
	}

	var (
		similarNames   []string
		containerFound *container.Summary
	)

	// Find the exact container and collect similar names
	for _, container := range containers {
		containerName := strings.TrimPrefix(container.Names[0], "/")

		if strings.EqualFold(containerName, name) {
			containerFound = &container
			break
		}

		if strings.Contains(strings.ToLower(containerName), strings.ToLower(name)) {
			similarNames = append(similarNames, containerName)
		}
	}

	// If exact match found, toggle the container
	if containerFound != nil {
		fmt.Printf("\n🔄 Toggling container: %s\n", strings.TrimPrefix(containerFound.Names[0], "/"))
		if err := dockerClient.ToggleContainer(containerFound.ID); err != nil {
			log.Fatalf("❌ Failed to toggle container %s: %v", containerFound.Names[0], err)
		}
		return
	}

	// If no exact match, display similar names
	fmt.Printf("\n❌ Container with name '%s' not found.\n", name)
	if len(similarNames) > 0 {
		fmt.Println("🔍 Did you mean one of these?")
		for _, similarName := range similarNames {
			fmt.Printf("  - %s\n", similarName)
		}
	} else {
		fmt.Println("No similar containers found.")
	}

	os.Exit(1)
}

func ToggleContainerByID(id string) {
	if id == "" {
		log.Fatal("Id Cannot be empty")
	}

	fmt.Printf("\n🔄 Toggling container: %s\n", id)

	if err := dockerClient.ToggleContainer(id); err != nil {
		log.Fatalf("❌ Failed to toggle container %s: %v", id, err)
	}

	return
}
