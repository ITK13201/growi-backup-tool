package cmd

import (
	"github.com/itk13201/growi-backup-tool/domain"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "growi-backup-tool",
		Short: "GROWI backup tool",
		Long:  "GROWI backup tool",
	}
	rootCmd.PersistentFlags().BoolP("debug", "D", false, "Debug mode")

	cfg := domain.NewConfig()
	rootCmd.AddCommand(NewExpandCmd(cfg))

	return rootCmd
}
