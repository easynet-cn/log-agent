package repository

import "xorm.io/xorm"

func InitRepository(engine *xorm.Engine) {
	engine.Sync2(new(SeekInfo))
}
