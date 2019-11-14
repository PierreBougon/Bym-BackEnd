package websocket_communication

import "sync"

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
