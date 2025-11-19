package matcher

import (
	"container/list"
	"errors"
	"fmt"
)

type RobotUpdate struct {
	status  string
	robotId int
}

type RobotItem struct {
	robotId int
}

type RobotQueue struct {
	queue *list.List
	pos   map[int]*list.Element
}

func NewRobotQueue() *RobotQueue {
	return &RobotQueue{
		queue: list.New(),
		pos:   make(map[int]*list.Element),
	}
}

func (q *RobotQueue) Len() int {
	return q.queue.Len()
}

func (q *RobotQueue) Enqueue(r RobotItem) error {
	_, exists := q.pos[r.robotId]

	if exists {
		return fmt.Errorf("robot Id already queued %d", r.robotId)
	}
	el := q.queue.PushBack(r)
	q.pos[r.robotId] = el
	return nil
}

func (q *RobotQueue) Dequeue(rID int) error {
	el := q.pos[rID]

	if el == nil {
		return fmt.Errorf("robot of Id %d does not exist", rID)
	}

	q.queue.Remove(el)
	delete(q.pos, rID)
	return nil
}

func (q *RobotQueue) Pop() (*RobotItem, error) {
	if q.queue.Len() == 0 {
		return nil, errors.New("there are no robots available")
	}

	el := q.queue.Front()
	robotEl := el.Value.(RobotItem)
	err := q.Dequeue(robotEl.robotId)
	if err != nil {
		return nil, fmt.Errorf("smth went horribly wrong :%s", err.Error())
	}
	return &robotEl, nil
}
