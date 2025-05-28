package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[string]*websocket.Conn)
var mu sync.Mutex

type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/ws", handleConnection)
	go sendPings()

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// handleConnection handles incoming WebSocket connections
func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}

	// Set a read deadline to avoid indefinite blocking
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	// Get the username from the query parameters
	username := r.URL.Query().Get("username")
	if username == "" {
		log.Println("Username not provided")
		conn.Close()
		return
	}

	log.Printf("New connection established for user %s\n", username)
	mu.Lock()
	clients[username] = conn
	mu.Unlock()
	go readMessages(conn, username)
}

func readMessages(conn *websocket.Conn, username string) {
	for {
		var msg Message
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))

		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message from %s: %v\n", username, err)
			removeClient(username)
			return
		}

		log.Printf("Received message from %s: %s\n", msg.Username, msg.Content)
		if err := broadcastMessage(msg); err != nil {
			log.Printf("Error broadcasting message: %v\n", err)
		}
	}
}

func broadcastMessage(msg Message) error {
	mu.Lock()
	defer mu.Unlock()
	for username, conn := range clients {
		if username != msg.Username {
			if err := conn.WriteJSON(msg); err != nil {
				log.Printf("Error sending message to %s: %v\n", username, err)
				removeClient(username)
			}
		}
	}
	return nil
}

func removeClient(username string) {
	mu.Lock()
	defer mu.Unlock()
	delete(clients, username)
}
func sendPings() {
	for {
		time.Sleep(25 * time.Second)
		mu.Lock()
		for username, conn := range clients {
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Error sending ping to %s: %v\n", username, err)
				removeClient(username)
			}
		}
		mu.Unlock()
	}
}
