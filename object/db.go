package object

import (
	"fmt"

	"github.com/spf13/viper"
	"xorm.io/xorm"
)

var dbs map[string]*xorm.Engine

func InitDb(viper *viper.Viper) {
	dbConfigs := viper.GetStringMap("datasources")
	dbs = make(map[string]*xorm.Engine, len(dbConfigs))

	for k := range dbConfigs {
		if engine, err := xorm.NewEngine(viper.GetString(fmt.Sprintf("datasources.%s.driver", k)), viper.GetString(fmt.Sprintf("datasources.%s.url", k))); err != nil {
			panic("连接数据库失败, error=" + err.Error())
		} else {
			engine.SetMaxOpenConns(viper.GetInt(fmt.Sprintf("datasources.%s.maxOpenConns", k)))
			engine.SetMaxIdleConns(viper.GetInt(fmt.Sprintf("datasources.%s.maxIdleConns", k)))
			engine.ShowSQL(viper.GetBool(fmt.Sprintf("datasources.%s.showSQL", k)))

			dbs[k] = engine
		}
	}

}

func GetDatabases() map[string]*xorm.Engine {
	return dbs
}

func GetDB() *xorm.Engine {
	return dbs["log-agent"]
}
