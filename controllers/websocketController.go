package controllers

import (
	"fmt"
	"github.com/PierreBougon/Bym-BackEnd/websocket_communication"
	"net/http"
)

var ConnectWebSocket = func(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Trying to establish a connection")
	fmt.Println(r.Header.Get("Sec-WebSocket-Protocol")) //Grab the token from the header

	var wsPool *websocket_communication.WSPool
	wsPool = websocket_communication.GetWSPool()
	wsPool.CreateWebSocket(w, r)
}
