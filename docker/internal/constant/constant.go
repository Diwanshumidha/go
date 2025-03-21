package constant

type Commands struct {
	Use   string
	Short string
	Long  string
}

var Root = &Commands{
	Use:   "dock",
	Short: "Utility cli for docker packaged with features not in docker cli",
	Long: `Dock is a CLI tool that allows you to list Docker containers and start/stop them.

It can be used to list Docker containers and start/stop them.`,
}

var Toggle = &Commands{
	Use:   "toggle",
	Short: "Toggle a container",
	Long: `Toggle a container.

It can be used to toggle a container.`,
}

var New = &Commands{
	Use:   "new",
	Short: "Create a new container",
	Long: `Create a new container.

It can be used to create a new container.`,
}
