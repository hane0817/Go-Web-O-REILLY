package main

import (
	"time"

	"github.com/gorilla/websocket"
)

//クライアントはチャットを行なっている一人のユーザを表します

type client struct {
	//socketはこのクライアントのためのWebsocketです
	socket *websocket.Conn
	//sendはメッセージが送られるチャネルです
	send chan *message
	//roomはこのクライアントが参加しているチャットルームです
	room *room
	//userDataはユーザに関する情報を保持します
	userData map[string]interface{}
}

func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
