package cmd

import "github.com/spf13/cobra"

func Execute() error {
	var root = &cobra.Command{
		Short: "root",
		Long:  `root`,
	}

	// root.AddCommand(createMappingIndex)
	root.AddCommand(serviceCmd)
	return root.Execute()
}
