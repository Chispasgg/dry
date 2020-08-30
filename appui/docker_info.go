package appui

import (
	"bytes"
	"strconv"

	termui "github.com/gizak/termui"
	drydocker "github.com/moncho/dry/docker"
	drytermui "github.com/moncho/dry/ui/termui"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/go-units"

	"github.com/moncho/dry/ui"
	"github.com/olekukonko/tablewriter"
)

// DockerInformation provides Docker information.
type DockerInformation interface {
	DockerEnv() drydocker.Env
	Info() (types.Info, error)
	Version() (*types.Version, error)
}

//DockerInfo is a widget to show Docker info
type DockerInfo struct {
	drytermui.SizableBufferer
}

//NewDockerInfo creates a DockerInfo widget
func NewDockerInfo(docker DockerInformation) *DockerInfo {
	di := drytermui.NewParFromMarkupText(DryTheme, dockerInfo(docker))
	di.BorderTop = false
	di.BorderBottom = true
	di.BorderLeft = false
	di.BorderRight = false
	di.BorderFg = termui.Attribute(DryTheme.Footer)
	di.BorderBg = termui.Attribute(DryTheme.Bg)

	di.Height = 4
	di.Bg = termui.Attribute(DryTheme.Bg)
	di.TextBgColor = termui.Attribute(DryTheme.Bg)
	di.Display = false
	return &DockerInfo{di}
}

func dockerInfo(docker DockerInformation) string {
	version, _ := docker.Version()
	info, _ := docker.Info()

	swarmInfo := info.Swarm

	buffer := new(bytes.Buffer)

	rows := [][]string{
		{
			ui.Blue("Docker Host:"), ui.Yellow(docker.DockerEnv().DockerHost), "",
			ui.Blue("Docker Version:"), ui.Yellow(version.Version)},
		{
			ui.Blue("Cert Path:"), ui.Yellow(docker.DockerEnv().DockerCertPath), "",
			ui.Blue("APIVersion:"), ui.Yellow(version.APIVersion)},
		{
			ui.Blue("Verify Certificate:"), ui.Yellow(strconv.FormatBool(docker.DockerEnv().DockerTLSVerify)), "",
			ui.Blue("OS/Arch/Kernel:"), ui.Yellow(version.Os + "/" + version.Arch + "/" + version.KernelVersion)},
	}

	rows = addHostInfo(rows, info)
	rows = addSwarmInfo(rows, swarmInfo)
	table := tablewriter.NewWriter(buffer)
	table.SetBorder(false)
	table.SetColumnSeparator("")
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(rows)
	table.Render()
	return buffer.String()
}

func addSwarmInfo(rows [][]string, info swarm.Info) [][]string {
	firstRow := rows[0]
	secondRow := rows[1]
	thirdRow := rows[2]

	firstRow = append(firstRow,
		ui.Blue("Swarm:"),
		ui.Yellow(string(info.LocalNodeState)))
	if info.LocalNodeState != swarm.LocalNodeStateInactive {
		if info.ControlAvailable {
			secondRow = append(secondRow,
				ui.Blue("Node role:"),
				ui.Yellow(string(swarm.NodeRoleManager)))
		} else {
			secondRow = append(secondRow,
				ui.Blue("Node role:"),
				ui.Yellow(string(swarm.NodeRoleWorker)))
		}
		thirdRow = append(thirdRow,
			ui.Blue("Nodes:"),
			ui.Yellow(strconv.Itoa(info.Nodes)))
	}

	return [][]string{firstRow, secondRow, thirdRow}

}
func addHostInfo(rows [][]string, info types.Info) [][]string {
	firstRow := rows[0]
	secondRow := rows[1]
	thirdRow := rows[2]

	firstRow = append(firstRow,
		ui.Blue("Hostname:"),
		ui.Yellow(info.Name))
	secondRow = append(secondRow,
		ui.Blue("CPU:"),
		ui.Yellow(strconv.Itoa(info.NCPU)))
	thirdRow = append(thirdRow,
		ui.Blue("Memory:"),
		ui.Yellow(units.BytesSize(float64(info.MemTotal))))

	return [][]string{firstRow, secondRow, thirdRow}

}
