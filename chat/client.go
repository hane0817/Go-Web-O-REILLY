package main

import (
	"github.com/gorilla/websocket"
)

//クライアントはチャットを行なっている一人のユーザを表します

type client struct {
	//socketはこのクライアントのためのWebsocketです
	socket *websocket.Conn
	//sendはメッセージが送られるチャネルです
	send chan []byte
	//roomはこのクライアントが参加しているチャットルームです
	room *room
}
