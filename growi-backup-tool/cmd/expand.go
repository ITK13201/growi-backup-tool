package cmd

import (
	"github.com/itk13201/growi-backup-tool/controllers"
	"github.com/itk13201/growi-backup-tool/domain"
	"github.com/itk13201/growi-backup-tool/services"
	"github.com/spf13/cobra"
	"log"
)

func NewExpandCmd(cfg *domain.Config) *cobra.Command {
	expandCmd := &cobra.Command{
		Use:   "expand",
		Short: "Expand pages from json to markdown files.",
		Long:  "Expand pages from json to markdown files",
		Run: func(cmd *cobra.Command, args []string) {
			isDebug, _ := cmd.Flags().GetBool("debug")
			inputDir, _ := cmd.Flags().GetString("input-dir")
			outputDir, _ := cmd.Flags().GetString("output-dir")
			logger := services.NewLogger(isDebug)
			controller := controllers.NewExpandController(cfg, logger, inputDir, outputDir)
			controller.Run()
		},
	}
	expandCmd.PersistentFlags().StringP("input-dir", "i", "", "Input directory")
	expandCmd.PersistentFlags().StringP("output-dir", "o", "", "Output directory")
	if err := expandCmd.MarkPersistentFlagRequired("input-dir"); err != nil {
		log.Fatal(err)
	}
	if err := expandCmd.MarkPersistentFlagRequired("output-dir"); err != nil {
		log.Fatal(err)
	}
	return expandCmd
}
