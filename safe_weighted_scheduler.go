package loadbalance

import "sync"

// SafeWeightedScheduler 并发安全的WeightedScheduler
type SafeWeightedScheduler struct {
	nodes     []*Node
	curIndex  int // 当前索引值
	curWeight int // 当前权重
	weightGcd int // 权重值的最大公约数
	maxWeight int // 最大权重值
	sync.Mutex
}

// NewSafeWeightedScheduler 新建SafeWeightedScheduler
func NewSafeWeightedScheduler(nodes []*Node) *SafeWeightedScheduler {
	return &SafeWeightedScheduler{
		nodes:     nodes,
		weightGcd: gcdSlice(nodes),
		maxWeight: maxSlice(nodes),
	}
}

// Next 获取下一个节点，可并发调用
func (ws *SafeWeightedScheduler) Next() *Node {
	ws.Lock()
	nodes, i, w, n := ws.nodes, ws.curIndex, ws.curWeight, len(ws.nodes)
	for {
		if i == 0 {
			w -= ws.weightGcd
			if w <= 0 {
				w = ws.maxWeight
				if w <= 0 {
					ws.Unlock()
					return nil
				}
			}
		}
		node := nodes[i]
		if node.Weight >= w {
			ws.curIndex = increase(i, n)
			ws.curWeight = w
			ws.Unlock()
			return node
		}
		i = increase(i, n)
	}
}

var _ Scheduler = &SafeWeightedScheduler{}
