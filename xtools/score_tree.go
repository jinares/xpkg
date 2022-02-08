package xtools

import (
	"sync"
)

//ScoreTree 使用树状数组实现
type ScoreTree struct {
	MaxScore int
	Mem      []int64 // 树使用的数组
	sync.RWMutex
}

func NewScoreTree(maxScore int) *ScoreTree {
	return &ScoreTree{
		MaxScore: maxScore,
		Mem:      make([]int64, maxScore+10),
	}
}

func (tree *ScoreTree) getTreeValue(n int) int64 {
	if n > tree.MaxScore || n < 0 {
		return 0
	}
	return tree.Mem[n]
}

func (tree *ScoreTree) addTreeValue(n int, value int64) {
	if n > tree.MaxScore || n < 0 {
		return
	}
	tree.Mem[n] += value
}

func lowbit(x int) int {
	return x & (-x)
}

// 给对应的积分增加人数
func (tree *ScoreTree) Add(n int, value int64) {
	tree.Lock()
	defer tree.Unlock()
	if n == 0 {
		tree.addTreeValue(0, value)
		return
	}
	if n < 0 {
		return
	}
	for n <= tree.MaxScore {
		tree.addTreeValue(n, value)
		n += lowbit(n)
	}
}

// 获取[0, n]的总和
func (tree *ScoreTree) Query(n int) int64 {
	tree.RLock()
	defer tree.RUnlock()
	var sum int64 = 0
	for n > 0 {
		sum += tree.getTreeValue(n)
		n -= lowbit(n)
	}
	sum += tree.getTreeValue(0)
	return sum
}

// 相当于求区间(score, MaxScore]的和
func (tree *ScoreTree) getOrder(score int) int64 {
	r := tree.Query(tree.MaxScore) // 区间[0, MaxScore]的和
	l := tree.Query(score)         // 区间[0, score]的和
	return r - l
}

// 获取排名 调用方法GetOrder
func (tree *ScoreTree) GetOrder(score int) int64 {
	order := tree.getOrder(score)
	return order
}

// 获取参与排名总人数
func (tree *ScoreTree) GetTotal() int64 {
	return tree.Query(tree.MaxScore)
}
