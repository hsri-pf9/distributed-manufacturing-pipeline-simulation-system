package utils

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"io"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type SSEManager struct {
	clients map[chan string]bool
	mu      sync.Mutex
}

// NewSSEManager initializes the SSEManager
func NewSSEManager() *SSEManager {
	return &SSEManager{
		clients: make(map[chan string]bool),
	}
}

// RegisterClient starts an SSE connection and listens for events
func (s *SSEManager) RegisterClient(c *gin.Context) {
	clientChan := make(chan string, 10) // Buffered channel for events

	// Add client to the map
	s.mu.Lock()
	s.clients[clientChan] = true
	s.mu.Unlock()

	log.Println("[SSE] New client connected")

	// Set headers for SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	// Create a done channel to detect client disconnect
	done := make(chan struct{})

	// // Run a goroutine to detect client disconnection
	// go func() {
	// 	<-c.Request.Context().Done()
	// 	close(done) // Notify main function to remove the client
	// }()
	// ✅ Detect client disconnect
	go func() {
		select {
		case <-c.Request.Context().Done(): // If request context is canceled
			log.Println("[SSE] Client closed connection (context done)")
		case <-c.Writer.CloseNotify(): // If client closes the connection abruptly
			log.Println("[SSE] Client closed connection (CloseNotify)")
		}
		close(done) // Notify main function to remove the client
	}()

	// Stream updates to the client
	c.Stream(func(w io.Writer) bool {
		select {
		case msg, ok := <-clientChan:
			if !ok {
				return false
			}
			fmt.Fprintf(w, "data: %s\n\n", msg)
			w.(http.Flusher).Flush()
			return true
		case <-done: // Client disconnected
			return false
		}
	})

	// Cleanup when client disconnects
	s.mu.Lock()
	delete(s.clients, clientChan)
	s.mu.Unlock()
	close(clientChan)

	log.Println("[SSE] Client disconnected")
}

// // Broadcast updates to all clients
// func (s *SSEManager) BroadcastUpdate(message string) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	log.Printf("[SSE] Broadcasting update: %s\n", message)

// 	for clientChan := range s.clients {
// 		select {
// 		case clientChan <- message:
// 		default:
// 			log.Println("[SSE] Client channel full, dropping message")
// 		}
// 	}
// }

func (s *SSEManager) BroadcastUpdate(data interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Convert data to JSON
	message, err := json.Marshal(data)
	if err != nil {
		log.Printf("[SSE] Failed to encode JSON: %v\n", err)
		return
	}

	log.Printf("[SSE] Broadcasting update: %s\n", message)

	// Send to all connected clients
	for clientChan := range s.clients {
		select {
		case clientChan <- string(message):
		default:
			log.Println("[SSE] Client channel full, dropping message")
		}
	}
}
