package main

import (
	"3body/connections"
	"3body/forth"
	"3body/world"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hypebeast/go-osc/osc"
)

// Original structs
type ForthRequest struct {
	Input string `json:"input"`
}

type ForthResponse struct {
	Output []string    `json:"output"`
	Stack  forth.Stack `json:"stack"`
	Error  string      `json:"error,omitempty"`
}

// New structs for memory state serialization
type MemoryState struct {
	Objects []MemoryObject `json:"objects"`
}

type MemoryObject struct {
	Type        string `json:"type"` // "hed" or "nod"
	ID          string `json:"id"`
	X           int    `json:"x"`
	Y           int    `json:"y"`
	Message     string `json:"message,omitempty"`
	ConnectsToX *int   `json:"connectsToX"`         // Changed from connectsToX
	ConnectsToY *int   `json:"connectsToY"`         // Changed from connectsToY
	IsCurrent   bool   `json:"isCurrent,omitempty"` // Whether this nod is the current node for any head
}

var (
	globalStack  forth.Stack
	globalState  forth.State
	globalMemory *world.Memory2D
)

var allowedOrigins = map[string]bool{
	"http://localhost:5173": true,
	"http://localhost:5174": true,
}

func initializeForth() {
	// Initialize the world
	rows, cols := 20, 20
	clock := world.NewClock(100) // 100ms interval
	globalMemory = world.NewMemory2D(rows, cols)

	// Initialize Forth interpreter
	globalStack = forth.CreateStack()
	globalState = forth.CreateInitialState()

	client := osc.NewClient("localhost", 7001)

	// Add world dictionary words to forth state
	worldWords := world.DefineWorldDictionary(globalMemory, clock, client)
	for name, word := range worldWords {
		globalState.Dictionary[name] = word
	}
	clock.Start(globalMemory)
}

// New function to extract coordinates from node ID
func parseNodeID(id string) (x int, y int) {
	fmt.Sscanf(id, "%d,%d", &x, &y)
	return
}
func intPtr(i int) *int {
	return &i
}

// New function to get memory state
func getMemoryState() MemoryState {
	rows, cols := globalMemory.Dimensions()
	state := MemoryState{
		Objects: make([]MemoryObject, 0),
	}

	// Get current nodes for all heads
	currentNodes := make(map[string]bool)
	for _, hed := range globalMemory.GetHeads() {
		if hed.CurrentNod() != nil {
			currentNodes[hed.CurrentNod().ID()] = true
		}
	}

	// Add all nodes
	grid := globalMemory.GetGrid()
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			if nod := grid[y][x]; nod != nil {
				obj := MemoryObject{
					Type:      "nod",
					ID:        nod.ID(),
					X:         x,
					Y:         y,
					Message:   string(nod.Message()),
					IsCurrent: currentNodes[nod.ID()],
				}
				if next := nod.Next(); next != nil {
					nextX, nextY := parseNodeID(next.ID())
					obj.ConnectsToX = intPtr(nextX)
					obj.ConnectsToY = intPtr(nextY)
				}
				state.Objects = append(state.Objects, obj)
			}
		}
	}

	// Add all heads
	for _, hed := range globalMemory.GetHeads() {
		x, y := parseNodeID(hed.ID())
		obj := MemoryObject{
			Type: "hed",
			ID:   hed.ID(),
			X:    x,
			Y:    y,
		}
		if first := hed.FirstNod(); first != nil {
			firstX, firstY := parseNodeID(first.ID())
			obj.ConnectsToX = intPtr(firstX)
			obj.ConnectsToY = intPtr(firstY)
		}
		state.Objects = append(state.Objects, obj)
	}

	return state
}

func enableCors(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if allowedOrigins[origin] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
	}
}

// New SSE endpoint
func streamMemoryState(w http.ResponseWriter, r *http.Request) {
	enableCors(w, r)

	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create ticker for 10 times per second
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// Get notification when client disconnects
	done := r.Context().Done()

	// Send initial state immediately
	state := getMemoryState()
	if data, err := json.Marshal(state); err == nil {
		fmt.Fprintf(w, "data: %s\n\n", data)
		w.(http.Flusher).Flush()
	} else {
		log.Printf("Error marshaling initial state: %v", err)
		return
	}

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			state := getMemoryState()
			data, err := json.Marshal(state)
			if err != nil {
				log.Printf("Error marshaling state: %v", err)
				return
			}
			// Proper SSE format: "data: " prefix and double newline
			fmt.Fprintf(w, "data: %s\n\n", data)
			w.(http.Flusher).Flush()
		}
	}
}

func streamMessages(w http.ResponseWriter, r *http.Request) {
	log.Println("New client connected to message stream")

	// Set headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// CORS headers
	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = "http://localhost:5173"
	}
	w.Header().Set("Access-Control-Allow-Origin", origin)

	// Initialize message channel if needed
	if connections.HTTPMessageChannel == nil {
		connections.HTTPMessageChannel = make(chan connections.HTTPMessage)
	}

	// Watch for client disconnection
	done := r.Context().Done()

	for {
		select {
		case <-done:
			log.Println("Client disconnected")
			return
		case msg := <-connections.HTTPMessageChannel:
			data, _ := json.Marshal(connections.HTTPMessage{Type: msg.Type, Content: msg.Content})
			fmt.Fprintf(w, "data: %s\n\n", data)
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		}
	}
}

func evaluateForth(w http.ResponseWriter, r *http.Request) {
	enableCors(w, r)

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ForthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Interpret the input
	stack, state, output := forth.Interpret(req.Input, globalStack, globalState)

	// Update global state
	globalStack = stack
	globalState = state

	// Prepare response
	response := ForthResponse{
		Output: output,
		Stack:  stack,
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Send response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func main() {
	// Initialize Forth interpreter
	initializeForth()

	// Set up HTTP routes
	http.HandleFunc("/evaluate", evaluateForth)
	http.HandleFunc("/memory-stream", streamMemoryState)
	http.HandleFunc("/message-stream", streamMessages)

	// Start server
	port := ":8080"
	fmt.Printf("Starting Forth interpreter server on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
