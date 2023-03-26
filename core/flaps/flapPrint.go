package flaps

import (
	"fmt"
	"time"
)

const (
	// FlapPrintName FlapPrint 的名称
	FlapPrintName = "print"
)

// FlapPrint 打印日志的 Flaps, 实现 FlapAction 接口
type FlapPrint struct {
	// 日志内容
	Msg string `json:"msg"`
}

func (f *FlapPrint) PluginConfig() any {
	return map[string]any{"msg": f.Msg}
}

// Plugin 插件名
func (f *FlapPrint) Plugin() string {
	return FlapPrintName
}

// Condition 自身的启动条件
func (f *FlapPrint) Condition() bool {
	return true
}

// FromConfig 从配置生成 FlapAction
func (f *FlapPrint) FromConfig(config interface{}) error {
	// 从配置生成 FlapPrint
	conf := config.(map[string]any)
	// 设置日志内容
	f.Msg = conf["msg"].(string)
	// 返回 FlapPrint
	return nil
}

// Execute 执行 Flap
func (f *FlapPrint) Execute(retryAttempt int) (*time.Time, error) {
	// action: 打印日志
	fmt.Print(f.Msg)
	// 不用重试
	return nil, nil
}

var _ FlapAction = (*FlapPrint)(nil)

// init 初始化 FlapPrint
func init() {
	// 注册 FlapPrint
	RegisterFlapActionMaker("print", func(config interface{}) (FlapAction, error) {
		return &FlapPrint{}, nil
	})
}
