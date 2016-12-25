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
