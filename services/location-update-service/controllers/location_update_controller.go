package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/niemet0502/zapp/services/location-update-service/services"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections (adjust for security!)
	},
}

type LocationUpdateController struct {
	svc *services.LocationUpdateService
}

func CreateLocationUpdatHandler(svc *services.LocationUpdateService) *LocationUpdateController {
	return &LocationUpdateController{svc}
}

func (h *LocationUpdateController) WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// Read messages in a loop
	for {
		var update services.Update
		err := conn.ReadJSON(&update)
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}

		h.svc.UpdateLocation(update)

		// Echo the message back to client
		// if err := conn.WriteMessage(messageType, message); err != nil {
		// 	fmt.Println("Write error:", err)
		// 	break
		// }
	}
}
