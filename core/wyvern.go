package core

import (
	"context"
	"errors"
)

var (
	//ErrSoarNotFound - Soar 不存在
	ErrSoarNotFound = errors.New("soar not found")
)

// Wyvern 结构体表示 Wyvern 编排系统
type Wyvern struct {
	// Soar 清单, id => Soar
	Soars map[string]*Soar
}

// NewWyvern 创建一个 Wyvern
func NewWyvern() *Wyvern {
	return &Wyvern{
		Soars: map[string]*Soar{},
	}
}

// Run 运行指定名称的 Soar
func (w *Wyvern) Run(ctx context.Context, name string) error {
	// 获取指定名称的 Soar
	soar, ok := w.Soars[name]
	if !ok {
		return ErrSoarNotFound
	}
	// 运行 Soar
	soar.Soar(ctx)
	return nil
}

// LoadFromConfig 从 WyvernConfig 配置加载某个名字的 Soar, 并返回其 id
func (w *Wyvern) LoadFromConfig(conf *WyvernConfig, name string) (string, error) {
	// 遍历获取指定名称的 Soar 配置
	soarConf, ok := conf.GetSoarConfByName(name)
	if !ok {
		return "", ErrSoarNotFound
	}
	// 使用 NewSoar 方法从配置创建 Soar
	soar, err := NewSoar(soarConf)
	if err != nil {
		return "", err
	}
	// 将 Soar 加入到 Wyvern 的 Soar 清单中
	w.Soars[soar.id] = soar
	return soar.id, nil
}
