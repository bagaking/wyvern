package flaps

import "time"

// FlapAction 定义 Flap 执行动作的函数签名
type FlapAction interface {
	// Execute 执行 Flap
	Execute(retryAttempt int) (*time.Time, error)

	// FromConfig 从配置生成 FlapAction
	FromConfig(config any) error

	// Condition 自身的启动条件
	Condition() bool

	// Plugin 名称
	Plugin() string

	// PluginConfig 配置的复制
	PluginConfig() any
}
