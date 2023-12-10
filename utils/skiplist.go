package utils

import (
	"math/rand"
	"sync"
)

const maxLevel = 4

type SkipList struct {
	head   *Element
	rwLock sync.RWMutex
}

type Element struct {
	K string
	V string

	score uint64
	// 记录每一层的下一个节点
	levels []*Element
}

func (e Element) Height() int {
	return len(e.levels)
}

func (s *SkipList) Set(k, v string) (isExist bool) {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	e := NewElement(k, v, randLevel())
	pre := s.head

	preElements := make([]*Element, len(s.head.levels))

	for i := len(s.head.levels) - 1; i >= 0; i-- {
		// 初始化pre节点
		preElements[i] = pre

		for cur := pre.levels[i]; cur != nil; cur = pre.levels[i] {
			if comp := compare(e, cur); comp <= 0 {
				if comp == 0 {
					cur.V = e.V
					isExist = true
					return
				}

				break
			}

			pre = cur
			preElements[i] = pre
		}
	}

	for i := 0; i < e.Height(); i++ {
		e.levels[i] = preElements[i].levels[i]
		preElements[i].levels[i] = e
	}

	return false
}

func (s *SkipList) Get(k string) (e *Element, ok bool) {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	pre := s.head

	for i := len(s.head.levels) - 1; i >= 0; i-- {
		for cur := pre.levels[i]; cur != nil; {
			if k < cur.K {
				break
			}
			if cur.K == k {
				return cur, true
			}
			pre = cur
			cur = cur.levels[i]
		}
	}

	return nil, false
}

func (s *SkipList) Del(k string) {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	pre := s.head

	preElements := make([]*Element, len(s.head.levels))
	var cur *Element

	for i := len(s.head.levels) - 1; i >= 0; i-- {
		for cur = pre.levels[i]; cur != nil; {
			if k < cur.K {
				break
			}
			if cur.K == k { // key 相同则覆盖
				preElements[i] = pre
				break
			}
			pre = cur
			cur = cur.levels[i]
		}
	}

	for i := len(preElements); i >= 0; i-- {
		preElements[i].levels[i] = cur.levels[i]
		cur.levels[i] = nil
	}
}

func NewSkipList() *SkipList {
	return &SkipList{
		head:   NewElement("", "", maxLevel),
		rwLock: sync.RWMutex{},
	}
}

func NewElement(k, v string, levels int) *Element {
	return &Element{
		K:      k,
		V:      v,
		score:  calKeyScore(k),
		levels: make([]*Element, levels),
	}
}

func calKeyScore(k string) uint64 {
	var hash uint64
	l := len(k)

	if l > 8 {
		l = 8
	}

	for i := 0; i < l; i++ {
		shift := uint(64 - 8 - i*8)
		hash |= uint64(k[i]) << shift
	}

	return hash
}

func compare(e1, e2 *Element) int {
	if e1.score < e2.score {
		return -1
	}
	if e1.score > e2.score {
		return 1
	}

	if e1.K < e2.K {
		return -1
	}
	if e1.K > e2.K {
		return 1
	}
	return 0
}

func randLevel() int {
	for i := 1; i <= maxLevel; i++ {
		if rand.Intn(2) < 1 {
			return i
		}
	}
	return maxLevel
}
