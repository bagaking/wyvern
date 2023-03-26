package core

type ID = string

// FlapIDTable 组件存储的索引结构, 包含 Soars Flaps 等, 并被所有子实体引用
type FlapIDTable map[ID]*Flap

func (w FlapIDTable) ListAllFlapID() []ID {
	// 列出所有 FlapID
	flapIDs := make([]ID, 0, len(w))
	for flapID := range w {
		flapIDs = append(flapIDs, flapID)
	}
	return flapIDs
}

func (w FlapIDTable) GetFlap(flapID ID) *Flap {
	// 获取指定 ID 的 Flap
	flap, ok := w[flapID]
	if !ok {
		return nil
	}
	return flap
}

// IFlapIndex - 行为索引接口
type IFlapIndex interface {
	// GetFlap - 获取指定 ID 的 Flap
	GetFlap(id ID) *Flap
	// ListAllFlapID - 列出所有 FlapID
	ListAllFlapID() []ID
}

var _ IFlapIndex = (*FlapIDTable)(nil)
