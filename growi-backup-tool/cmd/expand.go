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
			pagesFilePath, _ := cmd.Flags().GetString("pages-file")
			revisionsFilePath, _ := cmd.Flags().GetString("revisions-file")
			outputDirPath, _ := cmd.Flags().GetString("output-dir")
			cliArgument := &domain.CLIArgumentExpand{
				PagesFilePath:     pagesFilePath,
				RevisionsFilePath: revisionsFilePath,
				OutputDirPath:     outputDirPath,
			}
			logger := services.NewLogger(isDebug)
			controller := controllers.NewExpandController(cfg, logger, cliArgument)
			controller.Run()
		},
	}
	expandCmd.PersistentFlags().StringP("pages-file", "p", "", "pages file path")
	expandCmd.PersistentFlags().StringP("revisions-file", "r", "", "revisions file path")
	expandCmd.PersistentFlags().StringP("output-dir", "o", "", "output directory")
	if err := expandCmd.MarkPersistentFlagRequired("pages-file"); err != nil {
		log.Fatal(err)
	}
	if err := expandCmd.MarkPersistentFlagRequired("revisions-file"); err != nil {
		log.Fatal(err)
	}
	if err := expandCmd.MarkPersistentFlagRequired("output-dir"); err != nil {
		log.Fatal(err)
	}
	return expandCmd
}
