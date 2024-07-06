package growi

import (
	"bufio"
	"encoding/json"
	"github.com/itk13201/growi-backup-tool/domain"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type GrowiUtil struct {
	cfg       *domain.Config
	logger    *logrus.Logger
	inputDir  string
	outputDir string
}

func NewGrowiUtil(cfg *domain.Config, logger *logrus.Logger, inputDir string, outputDir string) *GrowiUtil {
	return &GrowiUtil{
		cfg:       cfg,
		logger:    logger,
		inputDir:  inputDir,
		outputDir: outputDir,
	}
}

func (g *GrowiUtil) loadTextFile(path string) (*string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	s := string(b)
	return &s, nil
}

func (g *GrowiUtil) loadPages() map[string]*domain.GrowiPage {
	growiPagesMap := map[string]*domain.GrowiPage{}

	pagesFilePath := filepath.Join(g.inputDir, g.cfg.Expand.PagesFileName)

	pagesFile, err := os.Open(pagesFilePath)
	if err != nil {
		g.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"path":  pagesFilePath,
		}).Fatal("failed to open pages file")
	}
	defer func(pagesFile *os.File) {
		err = pagesFile.Close()
		if err != nil {
			g.logger.WithFields(logrus.Fields{
				"error": err.Error(),
				"path":  pagesFilePath,
			}).Fatal("failed to close pages file")
		}
	}(pagesFile)

	scanner := bufio.NewScanner(pagesFile)
	for scanner.Scan() {
		line := scanner.Text()
		var pageMap map[string]interface{}
		err = json.Unmarshal([]byte(line), &pageMap)
		if err != nil {
			g.logger.WithFields(logrus.Fields{
				"error": err.Error(),
				"line":  line,
			}).Fatal("failed to unmarshal pageMap")
		}
		pageID := pageMap["_id"].(map[string]interface{})["$oid"].(string)
		var latestRevisionID *string
		if val, ok := pageMap["revision"]; ok {
			tmp := val.(map[string]interface{})["$oid"].(string)
			latestRevisionID = &tmp
		}
		growiPage := domain.GrowiPage{
			ID:               pageID,
			Path:             pageMap["path"].(string),
			LatestRevisionID: latestRevisionID,
		}
		growiPagesMap[pageID] = &growiPage
	}
	return growiPagesMap
}

func (g *GrowiUtil) dumpSinglePage(page *domain.GrowiPage) error {
	if page.LatestRevision == nil {
		// create only dir
		dirPath := filepath.Join(g.outputDir, page.Path)
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			g.logger.WithFields(logrus.Fields{
				"error": err.Error(),
				"path":  page.Path,
			}).Error("failed to create directory")
			return err
		}
	} else {
		body := page.LatestRevision.Body
		if body == "" {
			// create only dir
			dirPath := filepath.Join(g.outputDir, page.Path)
			err := os.MkdirAll(dirPath, os.ModePerm)
			if err != nil {
				g.logger.WithFields(logrus.Fields{
					"error": err.Error(),
					"path":  page.Path,
				}).Error("failed to create directory")
				return err
			}
		} else {
			// create page (*.md)
			lastElement := filepath.Base(page.Path)
			var filePath string
			if lastElement == "/" {
				filePath = filepath.Join(g.outputDir, "root.md")
			} else {
				filePath = filepath.Join(g.outputDir, page.Path+".md")
			}
			dirPath := filepath.Dir(filePath)
			err := os.MkdirAll(dirPath, os.ModePerm)
			if err != nil {
				g.logger.WithFields(logrus.Fields{
					"error":   err.Error(),
					"dirPath": dirPath,
				}).Error("failed to create directory")
				return err
			}
			file, err := os.Create(filePath)
			if err != nil {
				g.logger.WithFields(logrus.Fields{
					"error":    err.Error(),
					"filePath": filePath,
				}).Error("failed to create file")
				return err
			}
			defer func(file *os.File) {
				err = file.Close()
				if err != nil {
					g.logger.WithFields(logrus.Fields{
						"error": err.Error(),
					}).Error("failed to close file")
				}
			}(file)
			_, err = file.WriteString(body)
			if err != nil {
				g.logger.WithFields(logrus.Fields{
					"error":        err.Error(),
					"body(max:16)": body[:16],
				}).Error("failed to write body")
				return err
			}
		}
	}
	return nil
}

func (g *GrowiUtil) dumpPages(pagesMap map[string]*domain.GrowiPage) {
	revisionsFilePath := filepath.Join(g.inputDir, g.cfg.Expand.RevisionsFileName)

	revisionsFile, err := os.Open(revisionsFilePath)
	if err != nil {
		g.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"path":  revisionsFile,
		}).Fatal("failed to open revisions file")
	}
	defer func(revisionsFile *os.File) {
		err = revisionsFile.Close()
		if err != nil {
			g.logger.WithFields(logrus.Fields{
				"error": err.Error(),
				"path":  revisionsFile,
			}).Fatal("failed to close revisions file")
		}
	}(revisionsFile)

	scanner := bufio.NewScanner(revisionsFile)
	for scanner.Scan() {
		line := scanner.Text()
		var revisionMap map[string]interface{}
		err = json.Unmarshal([]byte(line), &revisionMap)
		if err != nil {
			g.logger.WithFields(logrus.Fields{
				"error": err.Error(),
				"line":  line,
			}).Fatal("failed to unmarshal revisionMap")
		}
		pageID := revisionMap["pageId"].(map[string]interface{})["$oid"].(string)
		// if pageID matched, dump page
		if page, ok := pagesMap[pageID]; ok {
			growiLatestRevision := domain.GrowiLatestRevision{
				ID:     revisionMap["_id"].(map[string]interface{})["$oid"].(string),
				Body:   revisionMap["body"].(string),
				PageID: pageID,
			}
			page.LatestRevision = &growiLatestRevision

			// dump page
			err = g.dumpSinglePage(page)
			if err != nil {
				g.logger.WithFields(logrus.Fields{
					"error":    err.Error(),
					"pageID":   pageID,
					"pagePath": page.Path,
				}).Error("failed to dump page")
			}
		}
	}
}

func (g *GrowiUtil) ExpandPages() {
	pagesMap := g.loadPages()
	g.dumpPages(pagesMap)
}
