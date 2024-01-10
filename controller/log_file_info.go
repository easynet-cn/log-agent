package controller

import (
	"net/http"

	"github.com/easynet-cn/log-agent/object"
	"github.com/easynet-cn/log-agent/repository"
	"github.com/gin-gonic/gin"
)

type logFileInfoController struct{}

var LogFileInfoController = new(logFileInfoController)

func (c *logFileInfoController) FindAll(ctx *gin.Context) {
	if entities, err := repository.LogFileInfoRepository.FindAll(object.GetDB()); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, entities)
	}
}
