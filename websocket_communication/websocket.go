package websocket_communication

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WebSocket struct {
	conn     *websocket.Conn
	clientId uint
	send     chan []byte
}

func (wsPool *WSPool) CreateWebSocket(w http.ResponseWriter, r *http.Request) *WebSocket {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Unable to upgrade connexion to websocket")
		fmt.Println(err)
		return nil
	}
	client := r.Context().Value("user").(uint)
	ws := WebSocket{conn: conn, clientId: client, send: make(chan []byte, 256)}
	go ws.readService()
	go ws.writeService()
	wsPool.AddSocket(&ws)
	return &ws
}

func (ws *WebSocket) readService() {
	defer func() {
		GetWSPool().RemoveSocket(ws)
		ws.conn.Close()
	}()
	ws.conn.SetReadLimit(maxMessageSize)
	ws.conn.SetReadDeadline(time.Now().Add(pongWait))
	ws.conn.SetPongHandler(func(string) error { ws.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := ws.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		print(message)
		ws.send <- message
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		//ws.hub.broadcast <- message
	}
}

func (ws *WebSocket) writeService() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		ws.conn.Close()
	}()
	for {
		select {
		case message, ok := <-ws.send:
			ws.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				ws.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := ws.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(ws.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-ws.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			ws.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
