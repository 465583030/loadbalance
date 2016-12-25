# loadbalance
负载均衡



## WeightedScheduler

加权轮询算法

假设有以下3个节点：

```
       ▅
       ▅
       ▅        ▅
       ▅        ▅
       ▅   ▅   ▅
权重    5    1    3
索引值  0    1    2
```

`WeightedScheduler` 每次都会选择权重值最大的节点，对于上面的例子选择的顺序是 `0, 0, 0, 2, 0, 2, 0, 1, 2`。

- 对于权重值相差比较大的节点，`WeightedScheduler` 会重复选择同一个节点。


- `Next` 方法不会返回权重值为负的节点。




示例：

```go
s := NewWeightedScheduler(BuildNodes([]int{1, 2, 3}))
node := s.Next()
fmt.Println(node.Data.(int))
```

并发安全的实例：

```go
s := NewSafeWeightedScheduler(BuildNodes([]int{1, 2, 3}))
node := s.Next()
fmt.Println(node.Data.(int))
```






## NginxScheduler

`NginxScheduler` nginx 上使用的负载均衡算法。

每次调用 `Next` 方法返回的节点的权重都会减去所有节点的有效权重之和，因此对于权重值相差比较大的节点也能有较好的均衡效果。

- `Next` 每次都会遍历所有节点，因此性能上有所损失。




示例：

```go
s := NewNginxScheduler(BuildNodes([]int{1, 2, 3}))
node := s.Next()
fmt.Println(node.Data.(int))
```

并发安全的示例：

```go
s := NewSafeNginxScheduler(BuildNodes([]int{1, 2, 3}))
node := s.Next()
fmt.Println(node.Data.(int))
```




## 性能测试

```bash
BenchmarkNginxScheduler_Next-4      	100000000	        11.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkWeightedScheduler_Next-4   	300000000	         5.50 ns/op	       0 B/op	       0 allocs/op
```



## 可视化

通过 web 页面查看各负载均衡算法的选择效果 [点击](https://github.com/zhangyuchen0411/loadbalance_web)

