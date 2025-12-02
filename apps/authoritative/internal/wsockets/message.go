package wsockets

type Message struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

type StatusMessage struct{}
