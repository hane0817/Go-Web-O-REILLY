package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/hane0817/Go-Web-O-REILLY.git/trace"
	"github.com/stretchr/objx"
)

type room struct {
	//forwardは他のクライアントに転送するためのメッセージを保持するチャネルです
	forward chan *message
	//joinはチャットルームに参加しようとしているクライアントのためのチャネルです
	join chan *client
	//leaveはチャットルームから退出しようとしているクライアントのためのチャネルです
	leave chan *client
	//clientsには在室している全てのクライアントが保持されます
	clients map[*client]bool
	//tracerはチャットルーム上で行われた操作のログを受け取ります
	tracer trace.Tracer
}

// newRoomはすぐに利用できるチャットルームを生成して返します
func newRoom() *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			//参加
			r.clients[client] = true

			r.tracer.Trace("新しいクライアントが参加しました")

		case client := <-r.leave:
			//退出
			delete(r.clients, client)
			close(client.send)

			r.tracer.Trace("クライアントが退出しました")

		case msg := <-r.forward:
			r.tracer.Trace("メッセージを受信しました: ", msg.Message)
			//全てのクライアントにメッセージを転送
			for client := range r.clients {
				select {
				case client.send <- msg:
					//メッセージを送信
					r.tracer.Trace("--クライアントに送信されました")
				default:
					//送信に失敗
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace("--送信に失敗しました。クライアントをクリーンアップします")
				}
			}
		}

	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("クッキーの取得に失敗しました:", err)
		return
	}

	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}

	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
