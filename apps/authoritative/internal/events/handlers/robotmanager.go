package handlers

import "errors"

type RobotAssignment struct {
	robotID int
	orderID int
}

func RobotAssigned(data []byte) error {
	// this will send an update via tcp to get thje robot to do shiot
	return errors.New("still workin on")
}
