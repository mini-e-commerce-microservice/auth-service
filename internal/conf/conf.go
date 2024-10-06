package conf

import (
	"flag"
	"github.com/go-faker/faker/v4"
	"github.com/spf13/viper"
)

var conf *Config

func Init() {
	if flag.Lookup("test.v") != nil {
		fakeConf := Config{}
		err := faker.FakeData(&fakeConf)
		if err != nil {
			panic(err)
		}
		conf = &fakeConf
		return
	}

	listDir := []string{".", "../", "../../", "../../../", "../../../../", "../../../../../", "../../../../../"}

	for _, dir := range listDir {
		viper.SetConfigName("env")
		viper.SetConfigType("json")
		viper.AddConfigPath(dir)
		err := viper.ReadInConfig()
		if err == nil {
			viper.SetConfigName("env.override")
			err = viper.MergeInConfig()
			if err != nil {
				panic(err)
			}
			if err = viper.Unmarshal(&conf); err != nil {
				panic(err)
			}
			return
		}
	}

	panic("cannot load env")
}

func InitForTest() {
	fakeConf := Config{}
	err := faker.FakeData(&fakeConf)
	if err != nil {
		panic(err)
	}
	conf = &fakeConf
	return
}

func GetConfig() *Config {
	return conf
}
