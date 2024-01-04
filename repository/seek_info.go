package repository

import "xorm.io/xorm"

type SeekInfo struct {
	Id      string `xorm:"varchar(50) 'id' pk notnull default('') comment('ID')" json:"id"`
	Ip      string `xorm:"varchar(20) 'ip' notnull default('') comment('IP')" json:"ip"`
	Project string `xorm:"varchar(50) 'project' notnull default('') comment('项目')" json:"project"`
	LogFile string `xorm:"varchar(1000) 'log_file' notnull default('') comment('日志文件绝对路径')" json:"logFile"`
	Offset  int64  `xorm:"bigint 'offset' notnull default(0) comment('偏移量')" json:"offset"`
}

type seekInfoRepository struct{}

var SeekInfoRepository = new(seekInfoRepository)

func (r *seekInfoRepository) Create(engine *xorm.Engine, entity *SeekInfo) (*SeekInfo, error) {
	_, err := engine.Insert(entity)

	return entity, err
}

func (r *seekInfoRepository) Update(engine *xorm.Engine, entity *SeekInfo) (*SeekInfo, error) {
	_, err := engine.ID(entity.Id).Where("ip=? AND log_file=?", entity.Ip, entity.LogFile).Cols("offset").Update(entity)

	return entity, err
}

func (r *seekInfoRepository) FindAll(engine *xorm.Engine) ([]SeekInfo, error) {
	entities := make([]SeekInfo, 0)

	err := engine.Find(&entities)

	return entities, err
}

func (r *seekInfoRepository) FindByIpAndLogFile(engine *xorm.Engine, ip string, logFile string) (*SeekInfo, error) {
	entity := &SeekInfo{}

	_, err := engine.Where("ip=? AND log_file=?", ip, logFile).Get(entity)

	return entity, err
}
