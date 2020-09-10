package llrb

import "container/list"

type filterSides byte

const (
	filterSideAt    filterSides = 1
	filterSideBelow             = 2
	filterSideAbove             = 4
)

type filterFunc func(item Item) filterSides
type enqueueFunc func(queue *list.List, pos *Node)

type Iterator struct {
	walkQueue *list.List
	enqueueF  enqueueFunc
}

func newIterator(root *Node, enqueueF enqueueFunc) *Iterator {
	walkQueue := list.New()
	walkQueue.PushFront(root)

	return &Iterator{
		walkQueue: walkQueue,
		enqueueF:  enqueueF,
	}
}

func (it *Iterator) Read() Item {
	if it.walkQueue == nil {
		return nil
	}

	var pos *list.Element

	for {
		pos = it.walkQueue.Front()
		if pos == nil {
			it.walkQueue = nil
			return nil
		}
		it.walkQueue.Remove(pos)

		switch posV := pos.Value.(type) {
		case *Node:
			it.enqueueF(it.walkQueue, posV)
		case Item:
			return posV
		}
	}
}

func (t *LLRB) ascendF(filterF filterFunc) *Iterator {
	enqueueF := func(q *list.List, pos *Node) {
		filterSides := filterF(pos.Item)
		if pos.Right != nil && (filterSides&filterSideAbove) > 0 {
			q.PushFront(pos.Right)
		}
		if (filterSides & filterSideAt) > 0 {
			q.PushFront(pos.Item)
		}
		if pos.Left != nil && (filterSides&filterSideBelow) > 0 {
			q.PushFront(pos.Left)
		}
	}

	return newIterator(t.root, enqueueF)
}

func (t *LLRB) descendF(filterF filterFunc) *Iterator {
	enqueueF := func(q *list.List, pos *Node) {
		filterSides := filterF(pos.Item)
		if pos.Left != nil && (filterSides&filterSideBelow) > 0 {
			q.PushFront(pos.Left)
		}
		if (filterSides & filterSideAt) > 0 {
			q.PushFront(pos.Item)
		}
		if pos.Right != nil && (filterSides&filterSideAbove) > 0 {
			q.PushFront(pos.Right)
		}
	}

	return newIterator(t.root, enqueueF)
}

var filterSideAny = filterSideBelow | filterSideAt | filterSideAbove

func filterAny(item Item) filterSides {
	return filterSideAny
}

func filterAbove(pivot Item) filterFunc {
	return func(item Item) filterSides {
		if !less(pivot, item) {
			return filterSideAbove
		}
		return filterSideAny
	}
}

func filterAtOrAbove(pivot Item) filterFunc {
	return func(item Item) filterSides {
		if less(item, pivot) {
			return filterSideAbove
		}
		return filterSideAny
	}
}

func filterBelow(pivot Item) filterFunc {
	return func(item Item) filterSides {
		if !less(item, pivot) {
			return filterSideBelow
		}
		return filterSideAny
	}
}

func filterAtOrBelow(pivot Item) filterFunc {
	return func(item Item) filterSides {
		if less(pivot, item) {
			return filterSideBelow
		}
		return filterSideAny
	}
}

func (t *LLRB) Ascend() *Iterator  { return t.ascendF(filterAny) }
func (t *LLRB) Descend() *Iterator { return t.descendF(filterAny) }

func (t *LLRB) AscendAbove(pivot Item) *Iterator  { return t.ascendF(filterAbove(pivot)) }
func (t *LLRB) DescendAbove(pivot Item) *Iterator { return t.descendF(filterAbove(pivot)) }

func (t *LLRB) AscendAtOrAbove(pivot Item) *Iterator  { return t.ascendF(filterAtOrAbove(pivot)) }
func (t *LLRB) DescendAtOrAbove(pivot Item) *Iterator { return t.descendF(filterAtOrAbove(pivot)) }

func (t *LLRB) AscendBelow(pivot Item) *Iterator  { return t.ascendF(filterBelow(pivot)) }
func (t *LLRB) DescendBelow(pivot Item) *Iterator { return t.descendF(filterBelow(pivot)) }

func (t *LLRB) AscendAtOrBelow(pivot Item) *Iterator  { return t.ascendF(filterAtOrBelow(pivot)) }
func (t *LLRB) DescendAtOrBelow(pivot Item) *Iterator { return t.descendF(filterAtOrBelow(pivot)) }
