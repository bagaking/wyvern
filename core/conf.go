package core

import (
	"github.com/bagaking/wyvern/core/flaps"
	"gopkg.in/yaml.v3"
)

// WyvernConfig - Wyvern 的配置
type WyvernConfig struct {
	// Soar 清单
	Soars []SoarConfig `yaml:"soars" json:"soars"`
}

// SoarConfig - Soar 的配置
type SoarConfig struct {
	// Soar 名称
	Name string `yaml:"name" json:"name"`
	// Flap 配置, 以 Prev/Next 表示 Flap 之间的关系, 平铺在一维数组中配置
	Flaps []flaps.FlapConfig `yaml:"flaps" json:"flaps"`
}

// GetSoarConfByName - 从配置中获取指定名称的 Soar 配置
func (conf *WyvernConfig) GetSoarConfByName(name string) (SoarConfig, bool) {
	for _, soar := range conf.Soars {
		if soar.Name == name {
			return soar, true
		}
	}
	return SoarConfig{}, false
}

// LoadWyvernConfig - 从文本中加载 Wyvern 配置
func LoadWyvernConfig(strConf string) (*WyvernConfig, error) {
	// 以 yml 格式解析 strConf
	conf := WyvernConfig{}
	err := yaml.Unmarshal([]byte(strConf), &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
