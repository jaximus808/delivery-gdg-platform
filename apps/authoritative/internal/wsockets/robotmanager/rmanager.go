package robotmanager

import (
	"log"
	"net/http"

	"github.com/jaximus808/delivery-gdg-platform/main/apps/authoritative/internal/wsockets"
)

func StartRobotManager() {
	hub := wsockets.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsockets.HandleWebSocket(hub, w, r)
	})

	addr := ":8080"

	log.Printf("websocket server starting at %s", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("Listen and Serve Failed", err)
	}
}
