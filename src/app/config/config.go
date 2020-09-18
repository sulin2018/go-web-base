package config

import (
	"time"
	// "log"
	"io/ioutil"

	"github.com/JeremyLoy/config"
	logs "github.com/sirupsen/logrus"
	"github.com/sulin2018/go-web-base/src/utils"
	"gopkg.in/yaml.v2"
)

type AppConf struct {
	AppName         string        `yaml:"AppName"`
	AppRunMode      string        `yaml:"AppRunMode"`
	AppAddr         string        `yaml:"AppAddr"`
	AppPort         int           `yaml:"AppPort"`
	AppReadTimeout  time.Duration `yaml:"AppReadTimeout"`
	AppWriteTimeout time.Duration `yaml:"AppWriteTimeout"`
	AppSecret       string        `yaml:"AppSecret"`
	AppCorsOrigin   []string      `yaml:"AppCorsOrigin"`

	LogFilePath string `yaml:"LogFilePath"`

	DBType           string `yaml:"DBType"`
	DBHost           string `yaml:"DBHost"`
	DBUser           string `yaml:"DBUser"`
	DBPassword       string `yaml:"DBPassword"`
	DBDatabase       string `yaml:"DBDatabase"`
	UserBasePassword string `yaml:"UserBasePassword"`
	PageSize         uint   `yaml:"PageSize"`
}

var AppConfig AppConf

func InitConfig(configFile string) {
	logs.Trace("config file is", configFile)

	if utils.IsNotExist(configFile) {
		logs.Fatalln("config file not exist!")
	}

	// load config from yaml
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		logs.Fatalln(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &AppConfig)
	if err != nil {
		logs.Fatalln(err.Error())
	}

	// load config from env
	err = config.FromEnv().To(&AppConfig)
	if err != nil {
		logs.Fatalln(err)
	}

	AppConfig.AppReadTimeout = AppConfig.AppReadTimeout * time.Second
	AppConfig.AppWriteTimeout = AppConfig.AppWriteTimeout * time.Second

	logs.Trace(AppConfig.AppName)
	logs.Trace("init config complate")
}
