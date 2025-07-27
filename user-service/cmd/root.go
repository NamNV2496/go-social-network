package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var port int

func Execute() error {
	var root = &cobra.Command{
		Use:   "app",
		Short: "root",
		Long:  `root`,
		Run: func(cmd *cobra.Command, args []string) {
			os.Setenv("USER_URL", "0.0.0.0:"+fmt.Sprint(port))
			Invoke(startServer).Run()
		},
	}

	// Default port is 5610
	root.Flags().IntVar(&port, "port", 5610, "the server port")

	root.AddCommand(createMappingIndex)
	root.AddCommand(migrateDataCmd)

	return root.Execute()
}
