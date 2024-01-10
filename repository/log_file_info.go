package repository

import "xorm.io/xorm"

type LogFileInfo struct {
	Id                string `xorm:"varchar(50) 'id' pk notnull default('') comment('ID')" json:"id"`
	Ip                string `xorm:"varchar(20) 'ip' notnull default('') comment('IP')" json:"ip"`
	Project           string `xorm:"varchar(50) 'project' notnull default('') comment('项目')" json:"project"`
	LogFile           string `xorm:"varchar(1000) 'log_file' notnull default('') comment('日志文件绝对路径')" json:"logFile"`
	LogstashTimestamp string `xorm:"varchar(100) 'logstash_timestamp' notnull default('') comment('Logstash Timestamp')" json:"logstashTimestamp"`
}

type logFileInfoRepository struct{}

var LogFileInfoRepository = new(logFileInfoRepository)

func (r *logFileInfoRepository) Create(engine *xorm.Engine, entity *LogFileInfo) (*LogFileInfo, error) {
	_, err := engine.Insert(entity)

	return entity, err
}

func (r *logFileInfoRepository) Update(engine *xorm.Engine, entity *LogFileInfo) (*LogFileInfo, error) {
	_, err := engine.ID(entity.Id).Where("ip=? AND log_file=?", entity.Ip, entity.LogFile).Cols("logstash_timestamp").Update(entity)

	return entity, err
}

func (r *logFileInfoRepository) FindAll(engine *xorm.Engine) ([]LogFileInfo, error) {
	entities := make([]LogFileInfo, 0)

	err := engine.Find(&entities)

	return entities, err
}

func (r *logFileInfoRepository) FindByIpAndLogFile(engine *xorm.Engine, ip string, logFile string) (*LogFileInfo, error) {
	entity := &LogFileInfo{}

	_, err := engine.Where("ip=? AND log_file=?", ip, logFile).Get(entity)

	return entity, err
}
