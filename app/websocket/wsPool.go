package websocket

import (
	"github.com/PierreBougon/Bym-BackEnd/app/models"
	"sync"
)

var once sync.Once

type WSPool struct {
	sockArr []*WebSocket
}

//Unique instance
var poolInstance *WSPool

func GetWSPool() *WSPool {

	once.Do(func() {
		poolInstance = new(WSPool)
	})

	return poolInstance
}

func (wsPool *WSPool) AddSocket(socket *WebSocket) {
	wsPool.sockArr = append(wsPool.sockArr, socket)
}

func (wsPool *WSPool) RemoveSocket(socket *WebSocket) {
	i := 0
	for ; i < len(wsPool.sockArr) && socket != wsPool.sockArr[i]; i++ {
	}
	if i >= len(wsPool.sockArr) {
		return
	}
	if i+1 >= len(wsPool.sockArr) {
		wsPool.sockArr = wsPool.sockArr[:i]
		return
	}
	wsPool.sockArr = append(wsPool.sockArr[:i], wsPool.sockArr[i+1:]...)
}

func (wsPool *WSPool) GetSocket(client uint) *WebSocket {
	i := 0
	for ; i < len(wsPool.sockArr) && client != wsPool.sockArr[i].clientId; i++ {
	}
	if i >= len(wsPool.sockArr) {
		return nil
	}
	return wsPool.sockArr[i]
}

func (wsPool *WSPool) BroadcastMessage(authorId uint, playlistId uint, message string) {
	followers := models.GetFollowers(playlistId)
	owner := models.GetPlaylistById(playlistId, nil).UserId

	if owner != authorId {
		wsPool.addMessageToQueue(owner, message)
	}
	for _, follower := range followers {
		wsPool.addMessageToQueue(follower.AccountId, message)
	}
}

func (wsPool *WSPool) addMessageToQueue(subscriber uint, message string) {
	ws := wsPool.GetSocket(subscriber)
	if ws != nil {
		ws.send <- []byte(message)
	}
}
