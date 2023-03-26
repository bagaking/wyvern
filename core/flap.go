package core

import (
	"context"
	"errors"
	"github.com/bagaking/wyvern/core/flaps"
	"time"
)

var (

	//	ErrFlapIsNotReady 表示 Flap 还未准备好, Tick 方法不会执行
	ErrFlapIsNotReady = errors.New("flap is not ready")
	// ErrFlapAlreadySuccess 表示 Flap 已经成功, Tick 方法不会再执行, Soar 方法不应该收到这个错误
	ErrFlapAlreadySuccess = errors.New("flap is already success")
	//	ErrFlapAlreadyFailed 表示 Flap 已经失败, Tick 方法不会再执行, 收到这个错误后, Soar 方法会直接退出并抛出异常
	ErrFlapAlreadyFailed = errors.New("flap is already failed")

	//	ErrFlapParentsAreNotAllFinished 表示 Flap 的父节点还未全部完成, Tick 方法不会执行
	ErrFlapParentsAreNotAllFinished = errors.New("flap parents are not all finished")

	//	ErrFlapWaitForAware 表示 Flap 还未到达下次执行时间, Tick 方法不会执行
	ErrFlapWaitForAware = errors.New("flap wait for next aware time")
)

// FlapStatus 状态
type FlapStatus int

const (
	FlapStateWait = iota
	FlapStateStated
	FlapStateInProgress
	FlapStatusErrorAndRetry
	FlapStateSuccess
	FlapStateFailed
)

// Flap 原子能力载体
type Flap struct {
	index IFlapIndex // index flap 在 wyvern 中的索引

	ConfName string // Flap 配置名
	ID       string // Flap 名称

	PrevFlaps         []ID             // 父节点
	NextFlaps         []ID             // 子节点
	State             FlapStatus       // Flap 状态，0 表示未完成，1 表示已完成
	Start             time.Time        // Flap 开始时间，延迟任务从这个时间开始
	NextAwakeTime     *time.Time       // Flap 重试时间
	AttemptRetryCount int              // 记录 Retry 次数
	Action            flaps.FlapAction // Flap 执行动作函数
}

// IsCompleted 判断 Flap 是否已经完成, 无论成功或失败都算完成
func (f *Flap) IsCompleted() bool {
	return f.State == FlapStateSuccess || f.State == FlapStateFailed
}

// IsReady 判断 Flap 是否满足启动条件
func (f *Flap) IsReady() bool {
	if !f.CheckAllParentsSuccess() {
		return false
	}

	// 判断自身的启动条件
	if !f.Action.Condition() {
		return false
	}

	return true
}

// NewFlap 从插件名和 FlapConfig 创建 Flap
func NewFlap(config flaps.FlapConfig, store Store) (*Flap, error) {
	// 通过配置名实例化 FlapAction
	action, err := flaps.MakeFlapAction(config.Plugin, config.PluginConfig)
	if err != nil {
		return nil, err
	}

	// 创建 Flap
	return &Flap{
		ConfName:          config.Name,
		ID:                store.MakeFlapID(),
		State:             FlapStateWait,
		Start:             time.Now(),
		NextAwakeTime:     nil,
		AttemptRetryCount: 0,
		Action:            action,
	}, nil
}

// HasPrevOfID 判断当前节点是否有指定父节点
func (f *Flap) HasPrevOfID(id string) bool {
	for _, pID := range f.PrevFlaps {
		if pID == id {
			return true
		}
	}
	return false
}

// HasNextOfID 判断当前节点是否有指定子节点
func (f *Flap) HasNextOfID(id string) bool {
	for _, nID := range f.NextFlaps {
		if nID == id {
			return true
		}
	}
	return false
}

// AddPrev 添加一个父节点
func (f *Flap) AddPrev(parent *Flap) {
	// 检查父节点是否已经在当前节点的前驱节点列表中
	for _, pID := range f.PrevFlaps {
		if pID == parent.ID {
			return
		}
	}

	// 将当前节点添加到父节点的后继节点列表中
	parent.NextFlaps = append(parent.NextFlaps, f.ID)
	// 将父节点添加到当前节点的前驱节点列表中
	f.PrevFlaps = append(f.PrevFlaps, parent.ID)
}

// AddNext 添加一个子节点
func (f *Flap) AddNext(child *Flap) {
	// 检查子节点是否已经在当前节点的后继节点列表中
	for _, childID := range f.NextFlaps {
		if childID == child.ID {
			return
		}
	}
	// 将当前节点添加到子节点的前驱节点列表中
	child.PrevFlaps = append(child.PrevFlaps, f.ID)
	// 将子节点添加到当前节点的后继节点列表中
	f.NextFlaps = append(f.NextFlaps, child.ID)
}

// CheckAllParentsSuccess 检查当前节点的所有前驱节点是否已经完成执行
func (f *Flap) CheckAllParentsSuccess() bool {
	if f.PrevFlaps == nil || len(f.PrevFlaps) == 0 {
		return true
	}
	// 遍历当前节点的所有前驱节点
	for _, parentID := range f.PrevFlaps {
		if parent := f.index.GetFlap(parentID); parent.State != FlapStateSuccess {
			return false
		}
	}
	// 所有前驱节点均已完成执行，返回 true
	return true
}

// UpdateStatus 更新当前节点的执行状态
func (f *Flap) UpdateStatus(status FlapStatus, nextAwakeTime *time.Time) FlapStatus {
	// 更新当前节点的执行状态
	switch f.State = status; status {
	case FlapStateStated:
		f.Start = *nextAwakeTime
	case FlapStatusErrorAndRetry:
		f.AttemptRetryCount++
		f.NextAwakeTime = nextAwakeTime
	case FlapStateInProgress:
		f.NextAwakeTime = nextAwakeTime
	}
	return f.State
}

// Tick 周期性执行 Flap, 该方法会被 Soar 方法调用
func (f *Flap) Tick(ctx context.Context) error {
	// 如果当前节点正在 wait 状态,且所有前驱节点均已完成执行,则将当前节点状态更新为 in progress
	if f.State == FlapStateWait {
		// 检查父节点是否全部完成
		if !f.CheckAllParentsSuccess() {
			// 父节点未全部完成, 直接返回. Flap 方法会收到 ErrFlapParentsAreNotAllFinished 错误, 并不做处理,继续执行下一个 Flap
			return ErrFlapParentsAreNotAllFinished
		}
		// 蓄势: parent 全部完成开始记录 State 时间, 该时间可以用于 condition 判断
		tNow := time.Now()
		f.UpdateStatus(FlapStateInProgress, &tNow)
	}

	if f.State == FlapStateSuccess {
		return ErrFlapAlreadySuccess
	} else if f.State == FlapStateFailed {
		return ErrFlapAlreadyFailed
	}

	if !f.IsReady() {
		return ErrFlapIsNotReady
	}

	if f.State == FlapStateStated {
		tNow := time.Now()
		f.UpdateStatus(FlapStateInProgress, &tNow)
	}

	if time.Now().Before(*f.NextAwakeTime) {
		return ErrFlapWaitForAware
	}

	// 执行动作
	nextTime, err := f.Action.Execute(f.AttemptRetryCount)
	if err != nil {
		// 出错并稍后重试
		if nextTime != nil {
			f.UpdateStatus(FlapStatusErrorAndRetry, nextTime)
		}
		// 出错并退出
		f.UpdateStatus(FlapStateFailed, nextTime)
	}
	// 成功
	f.UpdateStatus(FlapStateSuccess, nextTime)
	return nil
}
