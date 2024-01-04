package router

import (
	"fmt"
	"runtime"

	"net/http"

	"github.com/easynet-cn/log-agent/configuration"
	"github.com/easynet-cn/log-agent/controller"
	"github.com/easynet-cn/log-agent/log"

	"github.com/gin-contrib/gzip"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	server := gin.Default()

	server.Use(ginzap.Ginzap(log.Logger, configuration.Config.GetString("date-time-format"), false))
	server.Use(gzip.Gzip(gzip.DefaultCompression))
	server.Use(Recovery)

	server.GET("/system/stats", Stats)
	server.GET("/files/download", controller.FileController.Download)
	server.GET("/seek-infos", controller.SeekInfoController.FindAll)

	return server
}

func Recovery(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Logger.Error("系统内部错误", zap.Error(fmt.Errorf("%v", r)))

			c.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "message": fmt.Errorf("%v", r)})
		}
	}()
	c.Next()
}

func Stats(ctx *gin.Context) {
	stats := &runtime.MemStats{}

	runtime.ReadMemStats(stats)

	ctx.JSON(http.StatusOK, stats)
}
