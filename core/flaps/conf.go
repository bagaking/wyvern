package flaps

// FlapConfig - Flap 的配置
type FlapConfig struct {
	// Flap 名称
	Name string `yaml:"name" json:"name"`
	// Flap 的插件名
	Plugin string `yaml:"plugin" json:"plugin"`
	// Flap 的插件配置
	PluginConfig interface{} `yaml:"pluginConfig" json:"pluginConfig"`
	// Flap 的父节点
	PrevFlaps []string `yaml:"prevFlaps" json:"prevFlaps"`
	// Flap 的子节点
	NextFlaps []string `yaml:"nextFlaps" json:"nextFlaps"`
	// Flap 的启动条件
	Conditions []string `yaml:"conditions" json:"conditions"`
}
