package controllers

import (
	"github.com/itk13201/growi-backup-tool/domain"
	"github.com/itk13201/growi-backup-tool/internal/growi"
	"github.com/sirupsen/logrus"
)

type ExpandController struct {
	cfg       *domain.Config
	logger    *logrus.Logger
	inputDir  string
	outputDir string
}

func NewExpandController(cfg *domain.Config, logger *logrus.Logger, inputDir string, outputDir string) *ExpandController {
	return &ExpandController{
		cfg:       cfg,
		logger:    logger,
		inputDir:  inputDir,
		outputDir: outputDir,
	}
}

func (ec *ExpandController) Run() {
	ec.logger.Info("Started expanding pages...")

	growiUtil := growi.NewGrowiUtil(ec.cfg, ec.logger, ec.inputDir, ec.outputDir)
	growiUtil.ExpandPages()

	ec.logger.Info("Finished expanding pages.")
}
