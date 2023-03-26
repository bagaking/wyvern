// Package store 用于序列化和存储 soar 和 flap 的数据
package core

// Store 用于序列化和存储 soar 和 flap 的数据
// Store 不创建 soar 和 flap, 只负责把数据保存到持久化存储中和从持久化存储中加载数据
// 只能从配置创建 soar 和 flap
type Store interface {

	// Rebuild 重建索引, 从持久化存储中加载 soar 和 flap 的数据, 并创建索引
	// Store 接口中只有这一个方法会创建 soar 和 flap, 当计算发生迁移, 重新选举 leader 后, 会调用这个方法从持久化存储中恢复数据
	Rebuild(soarID ID) (*FlapIDTable, error)

	// SaveSoar 保存 soar 的数据
	SaveSoar(soar *Soar) error
	// SaveFlap 保存 flap 的数据
	SaveFlap(flap *Flap) error

	// LoadSoar 加载 soar 的数据, 只会根据 soar 的 ID 加载 soar 的数据, 不会创建 soar 和 flap
	LoadSoar(soar *Soar, id ID) error
	// LoadFlap 加载 flap 的数据, 只会根据 flap 的 ID 加载 flap 的数据, 不会创建 flap
	LoadFlap(index IFlapIndex, id ID) error

	// MakeSoarID 创建 Soar 的 ID
	MakeSoarID() ID
	// MakeFlapID 创建 Flap 的 ID
	MakeFlapID() ID
}

// Store 恢复流程
// 1. 从持久化存储中加载 soar 的数据, 其中包含 root flap 的 ID
// 2. 从持久化存储中加载 flap 的数据, 包含每个 flap 的状态和 flap 之间的关系
