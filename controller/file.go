package controller

import (
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"time"

	"github.com/easynet-cn/log-agent/configuration"
	"github.com/easynet-cn/log-agent/util"
	"github.com/gin-gonic/gin"
)

type fileController struct{}

var FileController = new(fileController)

func (c *fileController) Download(ctx *gin.Context) {
	project := ctx.Query("project")
	projectLogPath := configuration.Config.GetString(fmt.Sprintf("projects.%s.log.path", project))
	projectLogFile := configuration.Config.GetString(fmt.Sprintf("projects.%s.log.file", project))

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

		ctx.Header("Content-Type", "application/octet-stream")
		ctx.Header("Content-Disposition", "attachment; filename="+downloadFile)
		ctx.Header("Content-Transfer-Encoding", "binary")

		ctx.File(logfile)
	}
}
