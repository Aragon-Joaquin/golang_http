package ws

type WsServer struct {
	clients    map[*Client]bool
	Register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

func NewWebsocketServer() *WsServer {
	return &WsServer{
		clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

func (ws *WsServer) RunHandlers() {
	for {
		select {
		case client := <-ws.Register:
			ws.clients[client] = true

		case client := <-ws.unregister:
			if _, ok := ws.clients[client]; ok {
				delete(ws.clients, client)
				close(client.send)
			}

		case message := <-ws.broadcast:
			for client := range ws.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(ws.clients, client)

				}
			}
		}
	}

}
