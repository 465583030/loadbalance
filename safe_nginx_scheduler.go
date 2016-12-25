package loadbalance

import "sync"

// SafeNginxScheduler 并发安全的NginxScheduler
type SafeNginxScheduler struct {
	nodes          []*Node
	effWeightTotal int // 所有节点有效权重之和
	sync.Mutex
}

// NewSafeNginxScheduler 新建SafeNginxScheduler
func NewSafeNginxScheduler(nodes []*Node) *SafeNginxScheduler {
	effWeightTotal := 0
	for _, n := range nodes {
		if n.Weight < 0 {
			n.Weight = -n.Weight
			n.effective = n.Weight
		} else if n.Weight == 0 {
			n.effective = 1
		} else {
			n.effective = n.Weight
		}
		effWeightTotal += n.effective
	}
	return &SafeNginxScheduler{
		nodes:          nodes,
		effWeightTotal: effWeightTotal,
	}
}

// Next 返回下一个节点，可并发调用。
func (ns *SafeNginxScheduler) Next() *Node {
	var best *Node
	ns.Lock()
	for _, n := range ns.nodes {
		n.curWeight += n.effective // 每次检查都增加当前权重

		if n.effective < n.Weight {
			// 节点通了增加权重直到等于weight
			n.effective++
			ns.effWeightTotal++
		}
		if best == nil || n.curWeight > best.curWeight {
			// 选择当前权重最大的
			best = n
		}
	}

	if best == nil {
		ns.Unlock()
		return nil
	}
	best.curWeight -= ns.effWeightTotal // 被选中的减去有效权重之和, 这样下次该节点被选中的概率就小了
	ns.Unlock()
	return best
}

var _ Scheduler = &SafeNginxScheduler{}
