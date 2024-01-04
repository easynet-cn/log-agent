package object

import (
	"fmt"

	"path"
	"path/filepath"

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

			if logfile, err := filepath.Abs(path.Join(projectLogPath, projectLogFile)); err != nil {
				log.Logger.Error("获取日志文件绝对路径发生异常", zap.String("logfile", logfile), zap.Error(err))
			} else {
				offset := int64(0)
				ip := util.LocalIp()

				var seekInfoEntity *repository.SeekInfo

				if entity, err := repository.SeekInfoRepository.FindByIpAndLogFile(db, ip, logfile); err != nil {
					log.Logger.Error("查询日志文件偏移量发生异常", zap.String("logfile", logfile), zap.Error(err))
				} else if entity.Id == "" {
					seekInfoEntity = &repository.SeekInfo{Id: uuid.NewString(), Ip: ip, Project: project, LogFile: logfile, Offset: 0}

					if _, err := repository.SeekInfoRepository.Create(db, seekInfoEntity); err != nil {
						log.Logger.Error("创建日志文件偏移量发生异常", zap.String("logfile", logfile), zap.Error(err))
					}
				} else {
					seekInfoEntity = entity

					offset = entity.Offset
				}

				if t, err := tail.TailFile(logfile, tail.Config{Follow: true, Location: &tail.SeekInfo{Offset: offset, Whence: 0}}); err != nil {
					log.Logger.Error("读取日志文件发生异常", zap.String("logfile", logfile), zap.Error(err))
				} else {
					for line := range t.Lines {

						fmt.Println(line.Text)

						if _, err := repository.SeekInfoRepository.Update(db, &repository.SeekInfo{Id: seekInfoEntity.Id, Ip: ip, Project: project, LogFile: logfile, Offset: offset}); err != nil {
							log.Logger.Error("更新日志文件偏移量发生异常", zap.String("logfile", logfile), zap.Error(err))
						} else {
							seekInfoEntity.Offset = offset
							offset = offset + int64(len([]byte(line.Text)))
						}
					}
				}
			}
		}(project)
	}
}
