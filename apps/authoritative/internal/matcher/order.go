package matcher

import "container/heap"

type OrderItem struct {
	ownerId  int
	orderId  int
	orderNum int // this is the actual order number given for the day
}

type Item struct {
	Value    any
	Priority int
	Index    int
}

func CreateOrder(ownerId int, orderId int, orderNum int) *OrderItem {
	return &OrderItem{
		ownerId:  ownerId,
		orderId:  orderId,
		orderNum: orderNum,
	}
}

func (o *OrderItem) UpdateOrderNum(orderNum int) {
	o.orderNum = orderNum
}

type OrderQueue []*Item

func (pq OrderQueue) Len() int { return len(pq) }

func (pq OrderQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq OrderQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *OrderQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *OrderQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

type OrderPQ struct {
	h OrderQueue
}

func NewOrderPQ() *OrderPQ {
	pq := &OrderPQ{h: make(OrderQueue, 0)}
	heap.Init(&pq.h)
	return pq
}

func (pq *OrderPQ) Insert(orderItem *OrderItem) {
	item := &Item{Priority: orderItem.orderNum, Value: orderItem}
	heap.Push(&pq.h, item)
}

func (pq *OrderPQ) Pop() *OrderItem {
	if pq.Len() == 0 {
		return nil
	}
	return heap.Pop(&pq.h).(*Item).Value.(*OrderItem)
}

func (pq *OrderPQ) Len() int {
	return pq.h.Len()
}
