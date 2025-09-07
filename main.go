package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // allow all origins
}

var connMap = make(map[uint32]*websocket.Conn, 10)
var idCount = &atomic.Uint32{}

func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	idCount.Add(1)
	id := idCount.Load()
	connMap[id] = conn
	defer delete(connMap, id)

	fmt.Println("Client connected id =", id)

	// Echo messages back
	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err, "id =", id)
			break
		}
		fmt.Printf("Received: %s\n", msg)

		for connId, clientConn := range connMap {
			if connId == id {
				continue
			}

			clientConn.WriteMessage(mt, msg)
		}
	}

	fmt.Println("Client exiting", "id=", id)
}

func main() {
	http.HandleFunc("/ws", handleWS)
	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
