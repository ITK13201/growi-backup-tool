package controllers

import (
	"github.com/itk13201/growi-backup-tool/domain"
	"github.com/itk13201/growi-backup-tool/internal/growi"
	"github.com/sirupsen/logrus"
)

type ExpandController struct {
	cfg         *domain.Config
	logger      *logrus.Logger
	cliArgument *domain.CLIArgumentExpand
}

func NewExpandController(cfg *domain.Config, logger *logrus.Logger, cliArgument *domain.CLIArgumentExpand) *ExpandController {
	return &ExpandController{
		cfg:         cfg,
		logger:      logger,
		cliArgument: cliArgument,
	}
}

func (ec *ExpandController) Run() {
	ec.logger.Info("Started expanding pages...")

	growiUtil := growi.NewGrowiUtil(ec.cfg, ec.logger, ec.cliArgument)
	growiUtil.ExpandPages()

	ec.logger.Info("Finished expanding pages.")
}
