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
	"xorm.io/xorm"

	"github.com/easynet-cn/log-agent/configuration"
	"github.com/easynet-cn/log-agent/log"
	"github.com/easynet-cn/log-agent/repository"
	"github.com/easynet-cn/log-agent/util"
)

type LogFileWatchInfo struct {
	Project string `json:"project"`
	LogFile string `json:"logFile"`
	Ip      string `json:"ip"`
	Text    string `json:"text"`
}

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

				tailFile(db, startCarbon, project, logFile, ip, logFileInfoEntity.Id)
			}
		}(project)
	}
}

func tailFile(
	db *xorm.Engine,
	startCarbon carbon.Carbon,
	project string,
	logFile string,
	ip string,
	logFileInfoId string) {
	if t, err := tail.TailFile(logFile, tail.Config{Follow: true, MustExist: true}); err != nil {
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
					if bytes, err := json.Marshal(&LogFileWatchInfo{Project: project, LogFile: logFile, Ip: ip, Text: line.Text}); err != nil {
						log.Logger.Error("序列化日志文件信息发生异常", zap.String("logFile", logFile), zap.Error(err))
					} else {
						if _, err := repository.LogFileInfoRepository.Update(db, &repository.LogFileInfo{Id: logFileInfoId, Ip: ip, Project: project, LogFile: logFile, LogstashTimestamp: mMap["@timestamp"].(string)}); err != nil {
							log.Logger.Error("更新日志文件信息发生异常", zap.String("logFile", logFile), zap.Error(err))
						} else {
							producerName := configuration.Config.GetString(fmt.Sprintf("projects.%s.output.kafka.name", project))
							topic := configuration.Config.GetString(fmt.Sprintf("projects.%s.output.kafka.topic", project))

							Send(producerName, topic, bytes)
						}
					}
				}
			}
		}
	}
}
