package service

import (
	"dock/internal/dockerClient"
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/olekukonko/tablewriter"
)

// Max length for container names
const nameMaxLen = 18

func ListContainers() {
	containers, err := dockerClient.GetContainers(true)
	if err != nil {
		fmt.Println("❌ Error:", err)
		os.Exit(1)
	}

	if len(containers) == 0 {
		fmt.Println("❌ No containers found.")
		return
	}

	// Table setup
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "NAME", "STATUS", "IMAGE", "PORTS", "UPTIME"})

	// Modern style settings
	table.SetBorder(false)
	table.SetCenterSeparator("│")
	table.SetColumnSeparator("│")
	table.SetRowSeparator("─")
	table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderLine(true)

	// Colors for header
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiYellowColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiMagentaColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlueColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor},
	)

	// Iterate over the containers and format the data
	for _, container := range containers {
		name := trimName(container.Names[0][1:], nameMaxLen)
		status := colorStatus(container.State)
		ports := formatPorts(container.Ports)
		uptime := formatUptime(container.Status)

		table.Append([]string{
			container.ID[:12],
			name,
			status,
			trimName(container.Image, nameMaxLen),
			ports,
			uptime,
		})
	}

	table.Render()
}

func trimName(name string, maxLen int) string {
	if len(name) > maxLen {
		return name[:maxLen-3] + "..."
	}
	return name
}

func colorStatus(state string) string {
	const (
		green  = "\033[32m"
		red    = "\033[31m"
		yellow = "\033[33m"
		reset  = "\033[0m"
	)

	switch state {
	case "running":
		return green + "RUNNING" + reset
	case "exited":
		return red + "EXITED" + reset
	case "paused":
		return yellow + "PAUSED" + reset
	default:
		return state
	}
}

func formatPorts(ports []types.Port) string {
	if len(ports) == 0 {
		return "N/A"
	}

	var result []string
	for _, port := range ports {
		result = append(result, fmt.Sprintf("%s:%d->%d", port.IP, port.PublicPort, port.PrivatePort))
	}
	return strings.Join(result, ", ")
}

func formatUptime(status string) string {
	if !strings.HasPrefix(status, "Up ") {
		return "N/A"
	}
	uptime := strings.TrimPrefix(status, "Up ")
	return uptime
}
