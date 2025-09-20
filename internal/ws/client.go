package ws

import (
	"bytes"
	"time"

	"github.com/gorilla/websocket"
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

type Client struct {
	Conn      *websocket.Conn
	wsServer  *WsServer
	send      chan []byte
	ErrorChan chan error
}

func NewClient(conn *websocket.Conn, wsServer *WsServer) *Client {
	return &Client{
		Conn:      conn,
		wsServer:  wsServer,
		send:      make(chan []byte, 256),
		ErrorChan: make(chan error),
	}

}

func (c *Client) ReadPump() {
	defer func() {
		c.wsServer.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)              //maximum bytes of message
	c.Conn.SetReadDeadline(time.Now().Add(pongWait)) //reading timeout in unix
	c.Conn.SetPongHandler(func(string) error {
		return c.Conn.SetReadDeadline(time.Now().Add(pongWait)) // this was changed. it returned nil
	}) //server sends constantly pings to the clients, and immediately after clients need to respond with a "pong"

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.ErrorChan <- err
			}
			break
		}
		message = bytes.TrimSpace(bytes.ReplaceAll(message, []byte{'\n'}, []byte{' '}))
		c.wsServer.broadcast <- message
	}

}

func (c *Client) WritePump() {

}
