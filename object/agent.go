package object

import (
	"encoding/json"
	"fmt"

	"path"
	"path/filepath"

	"github.com/golang-module/carbon/v2"
	"github.com/google/uuid"
	"github.com/hpcloud/tail"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/easynet-cn/log-agent/log"
	"github.com/easynet-cn/log-agent/repository"
	"github.com/easynet-cn/log-agent/util"
)

func InitWatch(viper *viper.Viper) {
	projects := viper.GetStringMap("projects")

	db := GetDB()

	for project := range projects {
		go func(project string) {
			projectLogPath := viper.GetString(fmt.Sprintf("projects.%s.log.path", project))
			projectLogFile := viper.GetString(fmt.Sprintf("projects.%s.log.file", project))

			if logFile, err := filepath.Abs(path.Join(projectLogPath, projectLogFile)); err != nil {
				log.Logger.Error("获取日志文件绝对路径发生异常", zap.String("logFile", logFile), zap.Error(err))
			} else {
				var startCarbon carbon.Carbon
				var logFileInfoEntity *repository.LogFileInfo

				ip := util.LocalIp()

				if entity, err := repository.LogFileInfoRepository.FindByIpAndLogFile(db, ip, logFile); err != nil {
					log.Logger.Error("查询日志文件记录信息异常", zap.String("logFile", logFile), zap.Error(err))

					return
				} else if entity.Id == "" {
					if entity, err := repository.LogFileInfoRepository.Create(db, &repository.LogFileInfo{Id: uuid.NewString(), Ip: ip, Project: project, LogFile: logFile}); err != nil {
						log.Logger.Error("创建日志文件记录信息异常", zap.String("logFile", logFile), zap.Error(err))

						return
					} else {
						logFileInfoEntity = entity
					}
				} else {
					if entity.LogstashTimestamp != "" {
						startCarbon = carbon.Parse(entity.LogstashTimestamp)

						if startCarbon.Error != nil {
							return
						}
					}

					logFileInfoEntity = entity
				}

				if t, err := tail.TailFile(logFile, tail.Config{Follow: true}); err != nil {
					log.Logger.Error("读取日志文件发生异常", zap.String("logFile", logFile), zap.Error(err))
				} else {
					for line := range t.Lines {
						mMap := make(map[string]any)

						if err := json.Unmarshal([]byte(line.Text), &mMap); err != nil {
							log.Logger.Error("解析日志文件发生异常", zap.String("logFile", logFile), zap.String("text", line.Text), zap.Error(err))
						} else {
							t := carbon.Parse(mMap["@timestamp"].(string))

							if t.Error != nil {
								log.Logger.Error("解析日志文件LogstashTimestamp发生异常", zap.String("logFile", logFile), zap.String("text", line.Text), zap.Error(err))
							} else if startCarbon.IsInvalid() || (t.Gt(startCarbon) && startCarbon.IsValid()) {
								if _, err := repository.LogFileInfoRepository.Update(db, &repository.LogFileInfo{Id: logFileInfoEntity.Id, Ip: ip, Project: project, LogFile: logFile, LogstashTimestamp: mMap["@timestamp"].(string)}); err != nil {
									log.Logger.Error("更新日志文件信息发生异常", zap.String("logFile", logFile), zap.Error(err))
								} else {
									fmt.Println(line.Text)
								}
							}
						}
					}
				}
			}
		}(project)
	}
}
