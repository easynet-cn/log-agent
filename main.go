package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-module/carbon/v2"
	_ "github.com/mattn/go-sqlite3"

	"github.com/easynet-cn/log-agent/configuration"
	"github.com/easynet-cn/log-agent/log"
	"github.com/easynet-cn/log-agent/object"
	"github.com/easynet-cn/log-agent/repository"
	"github.com/easynet-cn/log-agent/router"

	"go.uber.org/zap"
)

func main() {
	configuration.InitConfiguration()
	log.InitLogger(configuration.Config)
	object.InitDb(configuration.Config)
	repository.InitRepository(object.GetDB())
	object.InitProducer(configuration.Config)

	time.LoadLocation(configuration.Config.GetString("time-zone"))
	carbon.SetTimezone(configuration.Config.GetString("time-zone"))

	go object.InitWatch(configuration.Config)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", configuration.Config.GetInt("server.port")),
		Handler:      router.NewRouter(),
		ReadTimeout:  0,
		WriteTimeout: 0,
	}

	log.Logger.Info("服务初始化成功")

	if err := server.ListenAndServe(); err != nil {
		log.Logger.Error("服务启动异常", zap.Error(err))
	}
}
