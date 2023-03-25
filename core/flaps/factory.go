package flaps

import "errors"

// PluginMaker 实例化方法接口
type PluginMaker func(config interface{}) (FlapAction, error)

var (
	// ErrPluginNotFound - 找不到 FlapAction 的实例化方法
	ErrPluginNotFound = errors.New("plugin not found")

	// pluginRegistry - PluginMaker 的注册表, key 为 plugin 名称, value 为 PluginMaker
	pluginRegistry = make(map[string]PluginMaker)
)

// RegisterFlapActionMaker 根据 plugin name 注册 FlapAction 实例化方法
func RegisterFlapActionMaker(name string, maker PluginMaker) {
	// 注册 FlapAction 实例化方法
	pluginRegistry[name] = maker
}

// GetFlapActionMaker 根据 plugin name 获取 FlapAction 实例化方法
func GetFlapActionMaker(name string) PluginMaker {
	// 获取 FlapAction 实例化方法
	return pluginRegistry[name]
}

// MakeFlapAction 根据 plugin name 和配置生成 FlapAction
func MakeFlapAction(plugin string, pluginConfig interface{}) (FlapAction, error) {
	// 根据 plugin name 获取 FlapAction 实例化方法
	maker := GetFlapActionMaker(plugin)
	if maker == nil {
		// 找不到 FlapAction 实例化方法
		return nil, ErrPluginNotFound
	}
	// 根据配置生成 FlapAction
	a, err := maker(pluginConfig)
	if err != nil {
		// 生成 FlapAction 失败
		return nil, err
	}
	// 加载 FlapAction 的配置
	err = a.FromConfig(pluginConfig)
	if err != nil {
		// 加载配置失败
		return nil, err
	}
	// 返回 FlapAction
	return a, nil
}
