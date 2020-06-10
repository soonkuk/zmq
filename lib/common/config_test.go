package common

import (
	"testing"

	"gopkg.in/yaml.v3"
)

var conf = &Conf{
	Total:          100,
	Fail:           10,
	ReportDuration: 10,
}

func TestGetConf(t *testing.T) {
	var c Conf
	c.GetConf()
	if *conf != c {
		s1, _ := yaml.Marshal(&conf)
		s2, _ := yaml.Marshal(&c)
		t.Errorf("wanted: %s, got: %s", string(s1), string(s2))
	}
}
