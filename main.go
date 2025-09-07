package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
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
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		file, err := os.ReadFile("index.html")
		if err != nil {
			http.Error(writer, "Unable to serve index.html"+err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = writer.Write(file)
		if err != nil {
			log.Println(err)
			return
		}
	})

	http.HandleFunc("/ws", handleWS)
	fmt.Println("Server started on :9933")
	http.ListenAndServe(":9933", nil)
}
