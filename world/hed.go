// world/hed.go
package world

import (
	"3body/forth"
	"fmt"
)

// Hed represents a head that moves through the nodes
type Hed struct {
	id         string
	first      *Nod // Points to first node in sequence
	current    *Nod // Current node in sequence
	last       *Nod // The last nod in a sequence, used for windowed nods, this is optional
	every      int  // How often to trigger
	bangs      int  // Count of bangs received
	stopped    bool // Whether head is stopped
	stack      forth.Stack
	forthState forth.State
	modifier   string // appended to the end of a message before execution
}

// NewHed creates a new Hed
func NewHed(id string, first *Nod, last *Nod, every int, state forth.State) (*Hed, error) {
	if id == "" {
		return nil, fmt.Errorf("hed id cannot be empty")
	}

	return &Hed{
		id:         id,
		first:      first,
		current:    first,
		every:      every,
		bangs:      0,
		stopped:    true,
		stack:      forth.CreateStack(),
		forthState: state,
	}, nil
}

// Bang processes a tick for this head
func (h *Hed) Bang() error {
	if h.stopped {
		return nil
	}

	h.bangs++
	if h.bangs%h.every == 0 {
		if h.current == nil {
			return fmt.Errorf("current node is nil")
		}

		// Process current node
		newStack, newState, _, err := h.current.Bang(h.stack, h.forthState, h.modifier)
		if err != nil {
			return fmt.Errorf("error processing node: %w", err)
		}

		h.stack = newStack
		h.forthState = newState

		// Move to next node or wrap around to first
		if h.current.Next() != nil {
			h.current = h.current.Next()
		} else if h.current == h.last {
			h.current = h.first
		} else {
			h.current = h.first
		}
	}
	return nil
}

// Start begins head movement
func (h *Hed) Start() {
	h.stopped = false
}

// Stop halts head movement
func (h *Hed) Stop() {
	h.stopped = true
}

// SetEvery updates the frequency
func (h *Hed) SetEvery(every int) {
	h.every = every
}

// SetModifier sets the modifier string
func (h *Hed) SetModifier(modifier string) {
	h.modifier = modifier
}

// ID returns the head's identifier
func (h *Hed) ID() string {
	return h.id
}

func (h *Hed) CurrentNod() *Nod {
	return h.current
}

func (h *Hed) FirstNod() *Nod {
	return h.first
}
