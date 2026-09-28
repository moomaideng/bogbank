package cli

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "template",
	Short: "Bogbank service template",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
