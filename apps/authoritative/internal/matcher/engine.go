package matcher

// this is the engine for our matching making service, in which when a robot becomes avaialbnle it will send an update, and every second this match maker will attempt to match a user and robot
// orders will come in from grpc request, and
import (
	"fmt"
	"time"
)

// matcher for orders and robots

type OrderRobotMatch struct {
	OrderID int
	RobotID string
}

type OrderRobotMatcher struct {
	orderIntake chan (*OrderItem)
	robotIntake chan (*RobotUpdate)
	orderQueue  *OrderPQ
	robotQueue  *RobotQueue
	orderCount  int64
}

func CreateOrderRobotMatcher() *OrderRobotMatcher {
	return &OrderRobotMatcher{
		orderIntake: make(chan (*OrderItem), 100),
		robotIntake: make(chan (*RobotUpdate), 100), // this should be a robot update
		orderQueue:  NewOrderPQ(),
		robotQueue:  NewRobotQueue(),
		orderCount:  0,
	}
}

func (orm *OrderRobotMatcher) SubmitOrder(o *OrderItem) {
	orm.orderIntake <- o
}

func (orm *OrderRobotMatcher) SubmitRobot(r *RobotUpdate) {
	orm.robotIntake <- r
}

func (orm *OrderRobotMatcher) attemptMatch(matchesChan chan (*OrderRobotMatch)) {
	if orm.orderQueue.Len() > 0 && orm.robotQueue.Len() > 0 { // we have at least one order and one robot available
		orderItem := orm.orderQueue.Pop()
		robotItem, err := orm.robotQueue.Pop()
		if err != nil {
			fmt.Println(err.Error())
		}

		matchesChan <- &OrderRobotMatch{
			OrderID: orderItem.orderId,
			RobotID: robotItem.robotID,
		}

		fmt.Printf("match created between orderId: %d, robotID %d\n", orderItem.orderId, robotItem.robotID)

	}
}

func (orm *OrderRobotMatcher) StartORM() chan *OrderRobotMatch {
	matchesQueue := make(chan (*OrderRobotMatch), 10)
	go orm.startEngine(matchesQueue)
	return matchesQueue
}

func (orm *OrderRobotMatcher) startEngine(matchesChan chan (*OrderRobotMatch)) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case orderReq := <-orm.orderIntake: // get order request
			orm.orderCount++
			orderReq.UpdateOrderNum(int(orm.orderCount))
			orm.orderQueue.Insert(orderReq) // put in heap

		case robotUpdate := <-orm.robotIntake:
			var err error

			if robotUpdate.status == "online" { // add to queue
				err = orm.robotQueue.Enqueue(RobotItem{
					robotID: robotUpdate.robotID,
				})
			} else {
				err = orm.robotQueue.Dequeue(robotUpdate.robotID) //
			}

			if err != nil {
				fmt.Println(err.Error())
			}

		case <-ticker.C:
			orm.attemptMatch(matchesChan)
		}
	}
}
