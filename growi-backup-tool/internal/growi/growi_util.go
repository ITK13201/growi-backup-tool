package growi

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/itk13201/growi-backup-tool/domain"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

type GrowiUtil struct {
	cfg               *domain.Config
	logger            *logrus.Logger
	cliArgumentExpand *domain.CLIArgumentExpand
}

func NewGrowiUtil(cfg *domain.Config, logger *logrus.Logger, cliArgumentExpand *domain.CLIArgumentExpand) *GrowiUtil {
	return &GrowiUtil{
		cfg:               cfg,
		logger:            logger,
		cliArgumentExpand: cliArgumentExpand,
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

	pagesFilePath := g.cliArgumentExpand.PagesFilePath

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
			IsDumped:         false,
		}
		growiPagesMap[pageID] = &growiPage
	}
	err = scanner.Err()
	if err != nil {
		g.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("failed to scan pages file")
	}

	return growiPagesMap
}

func (g *GrowiUtil) parseRevision(buf *bytes.Buffer) *domain.GrowiRevision {
	var revisionMap map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &revisionMap)
	if err != nil {
		g.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("failed to unmarshal revisionMap")
	}

	return &domain.GrowiRevision{
		ID:     revisionMap["_id"].(map[string]interface{})["$oid"].(string),
		Body:   revisionMap["body"].(string),
		PageID: revisionMap["pageId"].(map[string]interface{})["$oid"].(string),
	}
}

func (g *GrowiUtil) dumpSinglePage(page *domain.GrowiPage, revision *domain.GrowiRevision) error {
	outputDirPath := g.cliArgumentExpand.OutputDirPath

	body := revision.Body
	if body == "" {
		// create only dir
		dirPath := filepath.Join(outputDirPath, page.Path)
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
			filePath = filepath.Join(outputDirPath, "root.md")
		} else {
			filePath = filepath.Join(outputDirPath, page.Path+".md")
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
	return nil
}

func (g *GrowiUtil) dumpPages(pagesMap map[string]*domain.GrowiPage) {
	revisionsFilePath := g.cliArgumentExpand.RevisionsFilePath
	outputDirPath := g.cliArgumentExpand.OutputDirPath

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

	reader := bufio.NewReader(revisionsFile)
	for i := 0; ; i++ {
		firstBytes, isPrefix, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				g.logger.WithFields(logrus.Fields{
					"count": i,
					"error": err.Error(),
				}).Fatal("failed to read revisions file")
			}
		}

		buf := &bytes.Buffer{}
		_, err = buf.Write(firstBytes)
		if err != nil {
			g.logger.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Fatal("failed to load revisions file")
		}

		if isPrefix {
			// the line continues
			for {
				continuationBytes, isPrefix, err := reader.ReadLine()
				if err != nil {
					if err == io.EOF {
						break
					} else {
						g.logger.WithFields(logrus.Fields{
							"count": i,
							"error": err.Error(),
						}).Fatal("failed to read revisions file")
					}
				}

				_, err = buf.Write(continuationBytes)
				if err != nil {
					g.logger.WithFields(logrus.Fields{
						"count": i,
						"error": err.Error(),
					}).Fatal("failed to write revisions file")
				}

				if !isPrefix {
					// reached the end of the line
					break
				}
			}
		}

		revision := g.parseRevision(buf)
		if page, ok := pagesMap[revision.PageID]; ok {
			if *page.LatestRevisionID == revision.ID {
				// dump page
				g.logger.WithFields(logrus.Fields{
					"pageID":     page.ID,
					"revisionID": revision.ID,
					"path":       page.Path,
				}).Info("dumping page...")

				err = g.dumpSinglePage(page, revision)
				if err != nil {
					g.logger.WithFields(logrus.Fields{
						"error":      err.Error(),
						"pageID":     page.ID,
						"revisionID": revision.ID,
						"pagePath":   page.Path,
					}).Fatal("failed to dump page")
				}

				page.IsDumped = true

				g.logger.WithFields(logrus.Fields{
					"pageID":     page.ID,
					"revisionID": revision.ID,
					"path":       page.Path,
				}).Info("dumped page.")
			}
		}
	}

	// create dir of pages without revision
	for pageID, page := range pagesMap {
		if page.LatestRevisionID == nil && !page.IsDumped {
			dirPath := filepath.Join(outputDirPath, page.Path)
			err = os.MkdirAll(dirPath, os.ModePerm)
			if err != nil {
				g.logger.WithFields(logrus.Fields{
					"error":  err.Error(),
					"path":   page.Path,
					"pageID": pageID,
				}).Fatal("failed to create directory")
			} else {
				g.logger.WithFields(logrus.Fields{
					"pageID": pageID,
					"path":   page.Path,
				}).Info("created directory")
			}
		}
	}
}

func (g *GrowiUtil) ExpandPages() {
	pagesMap := g.loadPages()
	g.dumpPages(pagesMap)
}
