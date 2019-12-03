package websocket

import "fmt"

func NotifyPlaylistSubscribers(fromUser uint, playlist_id uint, message string) {
	wsPool := GetWSPool()
	wsPool.BroadcastMessage(fromUser, playlist_id, message)
}

func PlaylistNeedRefresh(playlist_id uint, author_id uint) string {
	return fmt.Sprintf("Playlist %d updated by %d", playlist_id, author_id)
}

func PlaylistDeleted(playlist_id uint, author_id uint) string {
	return fmt.Sprintf("Playlist %d deleted by %d", playlist_id, author_id)
}
