// world/worldDictionary_helpers.go
package world

import (
	"fmt"
	"lofl/forth"
)

// Helper functions for stack manipulation
func popString(stack forth.Stack) (string, forth.Stack, error) {
	if len(stack) == 0 {
		return "", stack, fmt.Errorf("stack underflow")
	}

	val := stack[len(stack)-1]
	stack = stack[:len(stack)-1]

	str, ok := val.(string)
	if !ok {
		return "", stack, fmt.Errorf("expected string, got %T", val)
	}

	return str, stack, nil
}

// Helper functions for stack manipulation
func popNumber(stack forth.Stack) (float64, forth.Stack, error) {
	if len(stack) == 0 {
		return 0, stack, fmt.Errorf("stack underflow")
	}

	val := stack[len(stack)-1]
	newStack := stack[:len(stack)-1]

	switch v := val.(type) {
	case float64:
		return v, newStack, nil
	case int:
		return float64(v), newStack, nil
	default:
		return 0, stack, fmt.Errorf("value is not a number: %v", val)
	}
}

// func getCoordinates(stack forth.Stack) (x, y int, newStack forth.Stack, err error) {
// 	if len(stack) < 2 {
// 		return 0, 0, stack, fmt.Errorf("stack underflow")
// 	}

// 	xVal, stack, err := popNumber(stack)
// 	yVal, stack, err := popNumber(stack)

// 	return int(xVal), int(yVal), stack, nil
// }

// Helper function to get array from stack
func getArray(stack forth.Stack) ([]interface{}, forth.Stack, error) {
	var arr []interface{}
	elementStack := stack[1:]

	newStack := make(forth.Stack, len(elementStack))
	copy(newStack, stack)

	for len(newStack) > 0 {
		item := newStack[len(newStack)-1]
		newStack = newStack[:len(newStack)-1]

		str, ok := item.(string)
		if ok && str == "[" {
			// Found start of array
			return arr, newStack, nil
		}
		arr = append([]interface{}{item}, arr...)
	}

	return nil, stack, fmt.Errorf("no array start marker found")
}

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
