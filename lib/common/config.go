package common

import (
	"log"

	"github.com/spf13/viper"
)

type Conf struct {
	Total          int    `yaml:"total"`
	Fail           int    `yaml:"fail"`
	ReportDuration int    `yaml:"reportDuration"`
	DeviceType     string `yaml:"deviceType"`
}

func (c *Conf) GetConf() *Conf {
	viper.SetDefault("total", DefaultTotalClient)
	viper.SetDefault("fail", DefaultFailClient)
	viper.SetDefault("reportDuration", DefaultReportDuration)
	viper.SetDefault("deviceType", DefaultDeviceType)
	viper.SetConfigFile("./configs/config.yml")
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
	if viper.IsSet("reportDuration") {
		c.ReportDuration = viper.GetInt("reportDuration")
	}
	if viper.IsSet("deviceType") {
		c.DeviceType = viper.GetString("deviceType")
	}

	return c
}
