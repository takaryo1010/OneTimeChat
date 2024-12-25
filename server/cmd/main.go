package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

func main() {
	// ルートディレクトリ
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	// WebSocket接続
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()

	fmt.Println("Chat server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("error:", err)
	}
}
