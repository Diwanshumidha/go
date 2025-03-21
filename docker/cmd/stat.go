package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

func GetStats() {
	app := tview.NewApplication()
	table := tview.NewTable().
		SetBorders(true).SetFixed(1, 1)

	table.SetTitle(" Docker CPU Usage ").SetBorder(true)

	// Header
	table.SetCell(0, 0, tview.NewTableCell("CONTAINER ID").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	table.SetCell(0, 1, tview.NewTableCell("NAME").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	table.SetCell(0, 2, tview.NewTableCell("IMAGE").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	table.SetCell(0, 3, tview.NewTableCell("CPU USAGE").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))

	go func() {
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}

		ctx := context.Background()

		for {
			containers, err := cli.ContainerList(ctx, container.ListOptions{All: false})
			if err != nil {
				panic(err)
			}

			// Clear table rows before updating

			table.Clear()
			// Reset header
			table.SetCell(0, 0, tview.NewTableCell("CONTAINER ID").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
			table.SetCell(0, 1, tview.NewTableCell("NAME").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
			table.SetCell(0, 2, tview.NewTableCell("IMAGE").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
			table.SetCell(0, 3, tview.NewTableCell("CPU USAGE").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
			table.SetCell(0, 4, tview.NewTableCell("MEM USAGE").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
			table.SetCell(0, 5, tview.NewTableCell("NETWORK").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
			table.SetCell(0, 6, tview.NewTableCell("Created").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))

			for i, container := range containers {
				stats, err := cli.ContainerStats(ctx, container.ID, false)
				if err != nil {
					continue
				}

				// Parse the CPU usage
				parsedStats := parseCPUUsage(stats)

				networkUsage := fmt.Sprintf("%s🔻 / %s 🔺", parsedStats.NetworkRx, parsedStats.NetworkTx)
				memoryUsage := fmt.Sprintf("%s / %s", parsedStats.Memory, parsedStats.MemoryLimit)

				created := time.Unix(container.Created, 0)

				table.SetCell(i+1, 0, tview.NewTableCell(container.ID[:12]).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignCenter))
				table.SetCell(i+1, 1, tview.NewTableCell(container.Names[0][1:]).SetTextColor(tcell.ColorLightCyan).SetAlign(tview.AlignCenter))
				table.SetCell(i+1, 2, tview.NewTableCell(container.Image).SetTextColor(tcell.ColorBlue).SetAlign(tview.AlignCenter))
				table.SetCell(i+1, 3, tview.NewTableCell(parsedStats.CPUPercent).SetTextColor(tcell.ColorOrange).SetAlign(tview.AlignCenter))
				table.SetCell(i+1, 4, tview.NewTableCell(memoryUsage).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
				table.SetCell(i+1, 5, tview.NewTableCell(networkUsage).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
				table.SetCell(i+1, 6, tview.NewTableCell(fmt.Sprintf("%s", created.Format("Jan 2 at 3:04 pm"))).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
			}

			app.Draw()
			time.Sleep(2 * time.Second) // Update every 2 seconds
		}
	}()

	if err := app.SetRoot(table, true).Run(); err != nil {
		panic(err)
	}
}

var statCmd = &cobra.Command{
	Use:   "stat",
	Short: "Shows the stats for all containers",
	Run: func(cmd *cobra.Command, args []string) {
		GetStats()
	},
}

func init() {
	rootCmd.AddCommand(statCmd)
}

type containerStats struct {
	CPUPercent       string
	Memory           string
	NetworkRx        string
	NetworkTx        string
	MemoryLimit      string
	MemoryEfficiency string
}

func parseCPUUsage(stats container.StatsResponseReader) containerStats {
	// Use the Docker API to read stats and parse CPU usage
	var data container.StatsResponse
	if err := json.NewDecoder(stats.Body).Decode(&data); err != nil {
		return containerStats{}
	}
	defer stats.Body.Close()

	// Calculate CPU usage
	cpuDelta := float64(data.CPUStats.CPUUsage.TotalUsage - data.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(data.CPUStats.SystemUsage - data.PreCPUStats.SystemUsage)
	cpuPercent := 0.0
	if systemDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(data.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}

	// Calculate Memory usage
	memoryUsageInBytes := float64(data.MemoryStats.Usage)
	memoryInMB := memoryUsageInBytes / 1024.0 / 1024.0

	memoryLimit := float64(data.MemoryStats.Limit) / 1024.0 / 1024.0
	memoryEfficiency := (memoryUsageInBytes / float64(data.MemoryStats.Limit)) * 100

	// ✅ Calculate Network usage
	var rxBytes, txBytes uint64
	if data.Networks != nil {
		for _, v := range data.Networks {
			rxBytes += v.RxBytes
			txBytes += v.TxBytes
		}
	}
	rxMB := float64(rxBytes) / 1024.0 / 1024.0
	txMB := float64(txBytes) / 1024.0 / 1024.0

	return containerStats{
		CPUPercent:       fmt.Sprintf("%.2f%%", cpuPercent),
		Memory:           fmt.Sprintf("%.2f MB", memoryInMB),
		MemoryLimit:      fmt.Sprintf("%.2f MB", memoryLimit),
		MemoryEfficiency: fmt.Sprintf("%.2f%%", memoryEfficiency),

		NetworkRx: fmt.Sprintf("%.2f MB", rxMB),
		NetworkTx: fmt.Sprintf("%.2f MB", txMB),
	}
}
