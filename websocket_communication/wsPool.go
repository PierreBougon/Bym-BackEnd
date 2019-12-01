package websocket_communication

import (
	"github.com/PierreBougon/Bym-BackEnd/models"
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
	for ; socket != wsPool.sockArr[i]; i++ {
	}
	wsPool.sockArr = append(wsPool.sockArr[:i], wsPool.sockArr[i+1:]...)
}

func (wsPool *WSPool) GetSocket(client uint) *WebSocket {
	i := 0
	for ; client != wsPool.sockArr[i].clientId; i++ {
	}
	return wsPool.sockArr[i]
}

func (wsPool *WSPool) BroadcastMessage(authorId uint, playlistId uint, message string) {
	followers := models.GetFollowers(playlistId)
	owner := models.GetPlaylistById(playlistId, nil).UserId

	if owner != authorId {
		wsPool.addMessageToQueue(authorId, owner, message)
	}
	for _, follower := range followers {
		wsPool.addMessageToQueue(authorId, follower.AccountId, message)
	}
}

func (wsPool *WSPool) addMessageToQueue(authorId uint, subscriber uint, message string) {
	if subscriber != authorId {
		ws := wsPool.GetSocket(subscriber)
		ws.send <- []byte(message)
	}
}
