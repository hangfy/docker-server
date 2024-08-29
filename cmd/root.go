package cmd

import "github.com/spf13/cobra"

func NewCmdRoot() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "docker-ctl",
		Short: "A CLI to manage Docker services",
		Long: `docker-ctl is a command-line tool for managing Docker services.
You can use it to start, stop, or restart a specific service.`,
	}
	rootCmd.AddCommand(NewCmdServer())
	return rootCmd
}
