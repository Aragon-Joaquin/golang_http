package main

import (
	"golang-http/internal/ws"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:   4096,
	WriteBufferSize:  4096,
	HandshakeTimeout: 2 * time.Second,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "http://localhost:8080"
	}}

func createConnection(wsServer *ws.WsServer, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	client := ws.NewClient(conn, wsServer)

	go client.WritePump()
	go client.ReadPump()

	wsServer.Register <- client

	//wait if there's an error
	wsErr := <-client.ErrorChan
	client.Conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseAbnormalClosure, wsErr.Error()))
}
