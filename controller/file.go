package controller

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/easynet-cn/log-agent/configuration"
	"github.com/easynet-cn/log-agent/util"
	"github.com/gin-gonic/gin"
)

type fileController struct{}

var FileController = new(fileController)

func (c *fileController) Files(ctx *gin.Context) {
	project := ctx.Query("project")
	projectLogPath := configuration.Config.GetString(fmt.Sprintf("projects.%s.log.path", project))

	if dirEntries, err := os.ReadDir(projectLogPath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
	} else {
		files := make([]string, 0, len(dirEntries))

		for _, dirEntry := range dirEntries {
			if !dirEntry.IsDir() {
				files = append(files, dirEntry.Name())
			}
		}

		ctx.JSON(http.StatusOK, files)
	}
}

func (c *fileController) Download(ctx *gin.Context) {
	project := ctx.Query("project")
	projectLogFile := ctx.Query("file")

	projectLogPath := configuration.Config.GetString(fmt.Sprintf("projects.%s.log.path", project))

	if projectLogFile == "" {
		projectLogFile = configuration.Config.GetString(fmt.Sprintf("projects.%s.log.file", project))
	}

	if logfile, err := filepath.Abs(path.Join(projectLogPath, projectLogFile)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
	} else {
		downloadFile := projectLogFile
		ext := filepath.Ext(projectLogFile)

		if ext == "" {
			downloadFile = fmt.Sprintf("%s-%s-%s%s", projectLogFile, util.LocalIp(), time.Now().Format("20060102150405"), ".log")
		} else {
			downloadFile = fmt.Sprintf("%s-%s-%s%s", projectLogFile[:len(projectLogFile)-len(ext)], util.LocalIp(), time.Now().Format("20060102150405"), ext)
		}

		c.download(ctx, logfile, downloadFile)
	}
}

func (c *fileController) download(ctx *gin.Context, file string, downloadFile string) {
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment; filename="+downloadFile)
	ctx.Header("Content-Transfer-Encoding", "binary")

	ctx.File(file)
}
