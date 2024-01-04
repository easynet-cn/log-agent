package configuration

import "github.com/spf13/viper"

var (
	Config *viper.Viper
)

func InitConfiguration() {
	Config = viper.New()

	Config.SetConfigName("application")
	Config.SetConfigType("yml")
	Config.AddConfigPath("./config")

	if err := Config.ReadInConfig(); err != nil {
		panic(err)
	}
}
