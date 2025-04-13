package cmd

import "github.com/spf13/cobra"

func Execute() error {
	var root = &cobra.Command{
		Short: "root",
		Long:  `root`,
		Run: func(cmd *cobra.Command, args []string) {
			Invoke(
				// StartGRPC,
				// StartREST,
				startServer,
			).Run()
		},
	}

	root.AddCommand(createMappingIndex)
	root.AddCommand(migrateDataCmd)
	return root.Execute()
}
