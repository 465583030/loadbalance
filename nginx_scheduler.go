package loadbalance

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

// NginxScheduler Nginx中使用的负载均衡算法，每次从所有节点中选择权重最高的节点，并将其权重值降低。
type NginxScheduler struct {
	nodes          []*Node
	effWeightTotal int // 所有节点有效权重之和
}

// NewNginxScheduler 新建NginxScheduler
func NewNginxScheduler(nodes []*Node) *NginxScheduler {
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
	return &NginxScheduler{
		nodes:          nodes,
		effWeightTotal: effWeightTotal,
	}
}

// Next 每次调用Next都会遍历所有节点。
// 选出的节点的当前权重会减去所有节点的有效权重之和。
// 对于节点间权重相差比较大的情况，NginxScheduler的选择效果比WeightedScheduler要好一些，更加均衡，但性能要差些。
func (ns *NginxScheduler) Next() *Node {
	var best *Node
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
		return nil
	}
	best.curWeight -= ns.effWeightTotal // 被选中的减去有效权重之和, 这样下次该节点被选中的概率就小了
	return best
}
