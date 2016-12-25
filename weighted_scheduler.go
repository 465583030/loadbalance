package loadbalance

// WeightedScheduler 加权轮询调度,每次从上次停止的地方选择权重最高的节点.
type WeightedScheduler struct {
	nodes     []*Node
	curIndex  int // 当前索引值
	curWeight int // 当前权重
	weightGcd int // 权重值的最大公约数
	maxWeight int // 最大权重值
}

// NewWeightedScheduler 新建WeightedScheduler
func NewWeightedScheduler(nodes []*Node) *WeightedScheduler {
	return &WeightedScheduler{
		nodes:     nodes,
		weightGcd: gcdSlice(nodes),
		maxWeight: maxSlice(nodes),
	}
}

// Next 选择下一个权重最高的节点.
//        ▅
//        ▅
//        ▅         ▅
//        ▅         ▅
//        ▅    ▅    ▅
// 权重    5    1    3
// 索引值  0    1    2
// 根据上面的权重，执行顺序应该是 0, 0, 0, 2, 0, 2, 0, 1, 2
// 形象的比喻就是先削顶层。
// 对于权重值很大的节点来说，Next 可能会重复选择该节点，此时建议选择NginxScheduler。
func (ws *WeightedScheduler) Next() *Node {
	nodes, i, w, n := ws.nodes, ws.curIndex, ws.curWeight, len(ws.nodes)
	for {
		if i == 0 {
			w -= ws.weightGcd
			if w <= 0 {
				w = ws.maxWeight
				if w <= 0 {
					return nil
				}
			}
		}
		node := nodes[i]
		if node.Weight >= w {
			ws.curIndex = increase(i, n)
			ws.curWeight = w
			return node
		}
		i = increase(i, n)
	}
}

func increase(i, n int) int {
	i++
	if i >= n {
		i = 0
	}
	return i
}

// 所有数的最大公约数
func gcdSlice(slice []*Node) int {
	n := len(slice)
	if n <= 0 {
		return 1
	}
	if n == 1 {
		return slice[0].Weight
	}
	g := slice[0].Weight
	for i := 1; i < n; i++ {
		g = gcd(g, slice[i].Weight)
	}
	return g
}

// 最大值
func maxSlice(slice []*Node) int {
	if len(slice) == 0 {
		return 0
	}
	m := slice[0].Weight
	for i := 1; i < len(slice); i++ {
		n := slice[i].Weight
		if n > m {
			m = n
		}
	}
	return m
}

// 最大公约数
func gcd(a, b int) int {
	if a <= 0 || b <= 0 {
		return 1
	}
	if a < b {
		a, b = b, a
	}
	for c := a % b; c > 0; {
		a, b = b, c
		c = a % b
	}
	return b
}
