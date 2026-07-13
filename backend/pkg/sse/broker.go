package sse

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Event represents a Server-Sent Event for task progress.
type Event struct {
	Type      string `json:"type"`
	TaskID    string `json:"taskId"`
	Status    string `json:"status"`
	Stage     string `json:"stage,omitempty"`
	Progress  int    `json:"progress"`
	Message   string `json:"message,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

// Broker manages SSE client subscriptions per task.
type Broker struct {
	mu      sync.RWMutex
	clients map[string]map[chan Event]struct{}
}

// NewBroker creates a new SSE broker.
func NewBroker() *Broker {
	return &Broker{
		clients: make(map[string]map[chan Event]struct{}),
	}
}

// Subscribe registers a new client channel for the given taskID.
func (b *Broker) Subscribe(taskID string) chan Event {
	ch := make(chan Event, 32)
	b.mu.Lock()
	if b.clients[taskID] == nil {
		b.clients[taskID] = make(map[chan Event]struct{})
	}
	b.clients[taskID][ch] = struct{}{}
	b.mu.Unlock()
	return ch
}

// Unsubscribe removes the client channel and closes it.
func (b *Broker) Unsubscribe(taskID string, ch chan Event) {
	b.mu.Lock()
	if clients, ok := b.clients[taskID]; ok {
		delete(clients, ch)
		if len(clients) == 0 {
			delete(b.clients, taskID)
		}
	}
	b.mu.Unlock()
	close(ch)
}

// Publish sends an event to all subscribers of the given taskID.
func (b *Broker) Publish(taskID string, event Event) {
	b.mu.RLock()
	clients := b.clients[taskID]
	b.mu.RUnlock()

	for ch := range clients {
		select {
		case ch <- event:
		default:
			// Client buffer full; skip
		}
	}
}

// StreamHandler returns a Gin handler that streams SSE events for a task.
func StreamHandler(broker *Broker) gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID := c.Param("id")
		ch := broker.Subscribe(taskID)
		defer broker.Unsubscribe(taskID, ch)

		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("X-Accel-Buffering", "no") // Disable Nginx buffering

		// Send heartbeat every 30 seconds
		heartbeat := time.NewTicker(30 * time.Second)
		defer heartbeat.Stop()

		c.Stream(func(w io.Writer) bool {
			select {
			case event, ok := <-ch:
				if !ok {
					return false
				}
				data, _ := json.Marshal(event)
				fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event.Type, data)
				return true
			case <-heartbeat.C:
				fmt.Fprintf(w, ": heartbeat\n\n")
				return true
			case <-c.Request.Context().Done():
				return false
			}
		})
	}
}

// ClientCount returns the number of active SSE connections for a task.
func (b *Broker) ClientCount(taskID string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.clients[taskID])
}

// TotalClients returns the total number of active SSE connections.
func (b *Broker) TotalClients() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	total := 0
	for _, clients := range b.clients {
		total += len(clients)
	}
	return total
}
