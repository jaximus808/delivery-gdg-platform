package robots

import "github.com/jaximus808/delivery-gdg-platform/main/apps/authoritative/internal/wsockets"

type RobotStage struct {
	robotID int
	client  *wsockets.Client
}
