package loadbalance

// Node 表示一个带权重的节点
type Node struct {
	Weight    int
	Data      interface{}
	effective int
	curWeight int
}

// Scheduler 调度器接口
type Scheduler interface {
	Next() *Node
}

// BuildNodes 一个辅助函数，根据权重列表生成节点列表，每个节点的Data表示其在权重列表中的索引。
func BuildNodes(weights []int) []*Node {
	nodes := make([]*Node, len(weights))
	for i, w := range weights {
		nodes[i] = &Node{
			Weight: w,
			Data:   i,
		}
	}
	return nodes
}
