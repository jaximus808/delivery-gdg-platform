package wsockets

type Message struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

type RobotUpdate struct {
	RobotID string `json:"robot_id"`
	Status  string `json:"status"`
}

type RobotMatch struct {
	RobotID string `json:"robot_id"`
	OrderID int    `json:"order_id"`
}
