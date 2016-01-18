package loadbalance

// 加权轮询调度
type WeightedRoundRobinScheduler struct {
	weights   []int // 权重
	curIndex  int   // 当前索引值
	curWeight int   // 当前权重
}

func NewWeightedRoundRobinScheduler(weights []int) *WeightedRoundRobinScheduler {
	return &WeightedRoundRobinScheduler{
		weights: weights,
	}
}

//        ▅
//        ▅
//        ▅         ▅
//        ▅         ▅
//        ▅    ▅    ▅
// 权重    5    1    3
// 索引值  0    1    2
// 根据上面的权重，执行顺序应该是 0, 0, 0, 2, 0, 2, 0, 1, 2
// 形象的比喻就是先削顶层。
func (this *WeightedRoundRobinScheduler) Next() (index int) {
	weights, i, w := this.weights, this.curIndex, this.curWeight
	for {
		if i == 0 {
			w -= this.gcdSlice(weights)
			if w <= 0 {
				w = this.maxSlice(weights)
				if w <= 0 {
					return -1
				}
			}
		}
		if weights[i] >= w {
			this.curIndex = (i + 1) % len(weights)
			this.curWeight = w
			return i
		}
		i = (i + 1) % len(weights)
	}
	return
}

// 所有数的最大公约数
func (this *WeightedRoundRobinScheduler) gcdSlice(slice []int) int {
	n := len(slice)
	if n <= 0 {
		return 1
	}
	if n == 1 {
		return slice[0]
	}
	gcd := slice[0]
	for i := 1; i < n; i++ {
		gcd = this.gcd(gcd, slice[i])
	}
	return gcd
}

// 最大公约数
func (this *WeightedRoundRobinScheduler) gcd(a, b int) int {
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

// 最大值
func (this *WeightedRoundRobinScheduler) maxSlice(slice []int) int {
	m := slice[0]
	for i := 1; i < len(slice); i++ {
		n := slice[i]
		if n > m {
			m = n
		}
	}
	return m
}
