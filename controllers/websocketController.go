package controllers

import (
	"github.com/PierreBougon/Bym-BackEnd/websocket_communication"
	"net/http"
)

var ConnectWebSocket = func(w http.ResponseWriter, r *http.Request) {
	print("Trying to establish a connection")
	var wsPool *websocket_communication.WSPool
	wsPool = websocket_communication.GetWSPool()
	wsPool.CreateWebSocket(w, r)
}
