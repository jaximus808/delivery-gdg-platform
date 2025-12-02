package matcher

import (
	"container/list"
	"errors"
	"fmt"
)

type RobotUpdate struct {
	status  string
	robotID string
}

func NewRobotUpdate(status string, robotID string) *RobotUpdate {
	return &RobotUpdate{
		status:  status,
		robotID: robotID,
	}
}

type RobotItem struct {
	robotID string
}

type RobotQueue struct {
	queue *list.List
	pos   map[string]*list.Element
}

func NewRobotQueue() *RobotQueue {
	return &RobotQueue{
		queue: list.New(),
		pos:   make(map[string]*list.Element),
	}
}

func (q *RobotQueue) Len() int {
	return q.queue.Len()
}

func (q *RobotQueue) Enqueue(r RobotItem) error {
	_, exists := q.pos[r.robotID]

	if exists {
		return fmt.Errorf("robot Id already queued %s", r.robotID)
	}
	el := q.queue.PushBack(r)
	q.pos[r.robotID] = el
	return nil
}

func (q *RobotQueue) Dequeue(rID string) error {
	el := q.pos[rID]

	if el == nil {
		return fmt.Errorf("robot of Id %s does not exist", rID)
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
	err := q.Dequeue(robotEl.robotID)
	if err != nil {
		return nil, fmt.Errorf("smth went horribly wrong :%s", err.Error())
	}
	return &robotEl, nil
}
