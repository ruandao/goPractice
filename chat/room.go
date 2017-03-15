package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	"github.com/ruandao/goPractice/trace"
)

type room struct {
	forward chan *message
	join chan *client
	leave chan *client
	clients map[*client]bool
	tracer trace.Tracer
}

func newRoom() *room {
	r := &room{
		forward:make(chan *message),
		join:	make(chan *client),
		leave:	make(chan *client),
		clients:make(map[*client]bool),
		tracer: trace.Off(),
	}
	return r
}

func (r *room) run() {
	for {
		select {
		case client := <- r.join:
			r.clients[client] = true
			r.tracer.Trace("New client joined")
		case client := <- r.leave:
			// 如果客户端断开连接，来不及发送leave事件怎么办
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("Client left")
		case msg := <- r.forward:
			r.tracer.Trace("Message received: ", msg.Message)
			for client := range r.clients{
				client.send <- msg
				r.tracer.Trace(" -- sent to client")
			}
		}
	}
}

const (
	socketBufferSize = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize:socketBufferSize, WriteBufferSize:socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println("ServeHTTP:", err)
		return
	}
	userData := make(map[string]interface{})
	cookie, err := req.Cookie("auth")
	if err != nil {
		log.Println("Failed to get auth cookie:", err)
		return
	}
	userData["name"] = cookie.Value
	userData["userid"] = cookie.Value

	cookie_avatar_url, err := req.Cookie("avatar_url")
	if err != nil {
		log.Println("Fail to get avatar_url from cookie: ", err)
	} else {
		userData["avatar_url"] = cookie_avatar_url.Value
	}

	client := &client{
		socket:	socket,
		send: 	make(chan *message, messageBufferSize),
		room:	r,
		userData:userData,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}