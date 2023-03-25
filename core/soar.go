package core

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Soar 结构体表示 Wyvern 中的原子能力
type Soar struct {
	RootFlaps []*Flap
	// 执行锁
	lock sync.Mutex
	// 执行次数
	count int
	// 创建时随机生成的独立uuid
	id string
}

// DFSUntil 遍历 DAG 对所有方法执行 action 直到指定条件满足, 或者遇到错误, 返回满足条件的 Flap
func (d *Soar) DFSUntil(ctx context.Context, action func(ctx context.Context, flap *Flap) (includeChildren bool, e error), condition func(flap *Flap) bool) (f *Flap, err error) {
	// 优雅退出并赋值 err
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	// 创建虚拟根节点, 用于遍历所有的根节点
	root := &Flap{
		Name: "__root__",
	}
	// 将所有的根节点加入到虚拟根节点的子节点中
	for _, flap := range d.RootFlaps {
		root.NextFlaps = append(root.NextFlaps, flap)
	}

	// 创建已访问列表
	visited := make(map[string]bool)
	// 加锁
	d.lock.Lock()
	defer d.lock.Unlock()

	// DFSUntil DFS 遍历 DAG
	var dfs func(ctx context.Context, flap *Flap) (*Flap, error)
	// 定义 dfs 函数
	dfs = func(ctx context.Context, flap *Flap) (*Flap, error) {
		// 如果已经在已访问列表中, 则说明存在环, 直接跳过
		if _, ok := visited[flap.Name]; ok {
			return nil, nil
		}
		// 如果这个节点已经满足要求了, 则直接返回
		if condition(flap) {
			return flap, nil
		}

		// 将自己加入到已访问列表中
		visited[flap.Name] = true

		// 除了虚拟节点, 其他执行 action
		if flap != root {
			if includeChildren, e := action(ctx, flap); e != nil {
				return nil, e
			} else if !includeChildren {
				// 如果 action 返回 false, 则不再遍历其子节点
				return nil, nil
			}
		}

		// 如果这个节点未完成, 则继续遍历其子节点
		for _, child := range flap.NextFlaps {
			// 如果这个子节点已经满足要求或者错误, 则直接返回
			if ff, e := dfs(ctx, child); e != nil || ff != nil {
				return ff, e
			}
		}

		// 如果这个节点的所有子节点都已经完成, 则返回 nil
		return nil, nil
	}

	// 进行 DFS 遍历
	return dfs(ctx, root)
}

// FindFlapByName 根据名称查找 Flap
func (d *Soar) FindFlapByName(name string) *Flap {
	// 使用 DFSUntil 查找
	f, _ := d.DFSUntil(context.Background(), func(ctx context.Context, flap *Flap) (bool, error) {
		return true, nil
	}, func(flap *Flap) bool {
		return flap.Name == name
	})
	return f
}

// Soar 启动一个协程，不断遍历 Flap DAG 并尝试执行最近未执行的项
func (d *Soar) Soar(ctx context.Context) {
	// 创建一个 goroutine
	go func(c context.Context) {
		// 发生异常时, 优雅退出
		defer func() {
			if r := recover(); r != nil {
				return
			}
		}()
		// 无限循环, 并记录执行次数和时间到context中
		for {
			//记录 id 执行次数 和 时间 到context中
			c = context.WithValue(c, "soar_id", d.id)
			c = context.WithValue(c, "soar_count", d.count)
			c = context.WithValue(c, "soar_time", time.Now())

			// 遍历 Flap DAG 并尝试执行最近未执行的项
			err := d.Flap(c)
			// 如果遇到错误, 则直接返回
			if err != nil {
				return
			}
			// 执行次数加一
			d.count++
			// 等待 500 毫秒
			time.Sleep(500 * time.Millisecond)
		}
	}(ctx)
}

// Flap 遍历 Flap DAG 并尝试执行最近未执行的项
// 每次遍历时，如果遇到一个未完成的 Flap，则执行该 Flap 的 Tick 方法， 否则将其子节点加入到 BFS 列表中
// 如果某个 Tick 方法返回的错误为 ErrFlapAlreadyFailed，则终止遍历
func (d *Soar) Flap(ctx context.Context) error {
	// 使用 DFSUntil 遍历 Flap DAG
	_, err := d.DFSUntil(ctx, func(ctx context.Context, flap *Flap) (bool, error) {
		// 如果这个 Flap 已经完成, 则直接跳过, 并继续遍历其子节点
		if flap.IsCompleted() {
			return true, nil
		}
		// 如果这个 Flap 未完成, 则执行 Tick 方法
		err := flap.Tick(ctx)
		// 如果 Tick 方法返回的错误为 ErrFlapAlreadyFailed，则终止遍历
		if err == ErrFlapAlreadyFailed {
			return false, err
		}
		// 不再遍历其子节点
		return false, nil
	}, func(flap *Flap) bool {
		// 不会被调用
		return false
	})

	// 返回错误
	return err
}

// NewSoar 从配置创建一个 Soar, 从配置文件中加载所有 Flap,并建立 Flap 之间的关系
func NewSoar(conf SoarConfig) (*Soar, error) {
	// 创建 Soar
	d := &Soar{
		RootFlaps: []*Flap{},
		lock:      sync.Mutex{},
		count:     0,
		id:        uuid.New().String(),
	}
	// 创建 Flap
	flaps := make(map[string]*Flap)
	for _, flapConf := range conf.Flaps {
		// 使用配置创建 Flap
		flap, err := NewFlap(flapConf)
		if err != nil {
			return nil, err
		}
		// 将 Flap 加入到 flaps 中
		flaps[flapConf.Name] = flap
	}
	// 建立 Flap 之间的关系
	for _, flapConf := range conf.Flaps {
		// 获取 Flap
		flap := flaps[flapConf.Name]
		// 设置 Flap 的 NextFlaps
		for _, nextFlapName := range flapConf.NextFlaps {
			flap.AddNext(flaps[nextFlapName])
		}
		// 设置 Flap 的 PrevFlaps
		for _, prevFlapName := range flapConf.PrevFlaps {
			flap.AddPrev(flaps[prevFlapName])
		}
	}
	// 将所有的根节点加入到 RootFlaps 中
	for _, flap := range flaps {
		// 找到入度为 0 的 Flap
		if len(flap.PrevFlaps) == 0 {
			d.RootFlaps = append(d.RootFlaps, flap)
		}
	}
	return d, nil
}
