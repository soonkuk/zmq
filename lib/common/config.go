package common

import (
	"log"

	"github.com/spf13/viper"
)

type Conf struct {
	Total          int    `yaml:"total"`
	Fail           int    `yaml:"fail"`
	Test           int    `yaml:"test"`
	Workers        int    `yaml:"workers"`
	ReportDuration int    `yaml:"reportDuration"`
	GroupID        string `yaml:"groupID"`
	DeviceType     string `yaml:"deviceType"`
	Port           string `yaml:"port"`
	EndPoint       string `yaml:"endpoint"`
}

func (c *Conf) GetConf(filePath string) *Conf {
	viper.SetDefault("total", DefaultTotalClient)
	viper.SetDefault("fail", DefaultFailClient)
	viper.SetDefault("test", DefaultTestClient)
	viper.SetDefault("reportDuration", DefaultReportDuration)
	viper.SetDefault("workers", DefaultWorkers)
	viper.SetDefault("groupID", DefaultGroupID)
	viper.SetDefault("deviceType", DefaultDeviceType)
	viper.SetDefault("port", DefaultPort)
	viper.SetDefault("endpoint", DefaultServerEndPoint)
	viper.SetConfigFile(filePath)
	viper.SetConfigType(DefaultConfigType)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		log.Printf("config file reading error : %s", err)
	}

	if viper.IsSet("total") {
		c.Total = viper.GetInt("total")
	}
	if viper.IsSet("fail") {
		c.Fail = viper.GetInt("fail")
	}
	if viper.IsSet("test") {
		c.Test = viper.GetInt("test")
	}
	if viper.IsSet("workers") {
		c.Workers = viper.GetInt("workers")
	}
	if viper.IsSet("reportDuration") {
		c.ReportDuration = viper.GetInt("reportDuration")
	}
	if viper.IsSet("deviceType") {
		c.DeviceType = viper.GetString("deviceType")
	}
	if viper.IsSet("groupID") {
		c.GroupID = viper.GetString("groupID")
	}
	if viper.IsSet("port") {
		c.Port = viper.GetString("port")
	}
	if viper.IsSet("endpoint") {
		c.EndPoint = viper.GetString("endpoint")
	}

	return c
}
