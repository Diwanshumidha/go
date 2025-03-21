package dockerClient

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type Options struct {
	Name     string
	User     string
	Password string
	DB       string
	Port     string
	Image    string
}

func CreatePostgresContainer(opts Options) (connectionString string, error error) {
	// Set default values
	if opts.Name == "" {
		opts.Name = "postgres-db"
	}
	if opts.User == "" {
		opts.User = "postgres"
	}
	if opts.Password == "" {
		opts.Password = "admin"
	}
	if opts.DB == "" {
		opts.DB = "postgres"
	}
	if opts.Port == "" {
		opts.Port = "5432"
	}
	if opts.Image == "" {
		opts.Image = "postgres:latest"
	}

	ctx := context.Background()

	// Create Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", fmt.Errorf("failed to create Docker client: %w", err)
	}

	// Check if container already exists
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return "", fmt.Errorf("failed to list containers: %w", err)
	}

	for _, c := range containers {
		if c.Names[0] == "/"+opts.Name {
			return "", fmt.Errorf("Container '%s' already exists.\n", opts.Name)
		}
	}

	// Define container configuration
	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: opts.Image,
			Env: []string{
				fmt.Sprintf("POSTGRES_USER=%s", opts.User),
				fmt.Sprintf("POSTGRES_PASSWORD=%s", opts.Password),
				fmt.Sprintf("POSTGRES_DB=%s", opts.DB),
			},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				nat.Port(fmt.Sprintf("%s/tcp", opts.Port)): {{HostIP: "0.0.0.0", HostPort: opts.Port}},
			},
			Mounts: []mount.Mount{},
			RestartPolicy: container.RestartPolicy{
				Name: "unless-stopped",
			},
		},
		&network.NetworkingConfig{},
		nil,
		opts.Name,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	// Start the container
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start container: %w", err)
	}

	fmt.Printf("PostgreSQL container '%s' started successfully.\n", opts.Name)

	cs := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", opts.User, opts.Password, "127.0.0.1", opts.Port, opts.DB)
	return cs, nil
}
