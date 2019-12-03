package controllers

import (
	"fmt"
	"github.com/PierreBougon/Bym-BackEnd/app/websocket"
	"net/http"
)

var ConnectWebSocket = func(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Trying to establish a connection")
	fmt.Println(r.Header.Get("Sec-WebSocket-Protocol")) //Grab the token from the header

	var wsPool *websocket.WSPool
	wsPool = websocket.GetWSPool()
	wsPool.CreateWebSocket(w, r)
}
