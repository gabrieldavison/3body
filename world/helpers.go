// world/worldDictionary_helpers.go
package world

import (
	"fmt"
)

// createSequence creates a sequence of nodes with messages
func createSequence(memory *Memory2D, x, y int, messages []string) error {
	// Reverse the messages array
	reversed := make([]string, len(messages))
	for i, j := 0, len(messages)-1; i < len(messages); i, j = i+1, j-1 {
		reversed[i] = messages[j]
	}

	// First, create all nodes and store them
	nodes := make([]*Nod, len(reversed))
	for i, message := range reversed {
		nod, err := NewNod(NodID(x+i, y), Message(message))
		if err != nil {
			return fmt.Errorf("error creating node: %w", err)
		}
		nodes[i] = nod
	}

	// Then, set up the connections between nodes
	for i := 0; i < len(nodes)-1; i++ {
		nodes[i].SetNext(nodes[i+1])
	}

	// Finally, add all nodes to memory
	for i, nod := range nodes {
		err := memory.AddNod(x+i, y, nod)
		if err != nil {
			return fmt.Errorf("error adding node: %w", err)
		}
	}

	return nil
}
