// nod.go
package world

import (
	"fmt"
	"lofl/forth"
	"strings"
)

type Message string

const (
	MessagePrint Message = "print"
	MessageSpeed Message = "speed"
)

type Nod struct {
	id      string
	message Message
	next    *Nod // Changed to pointer to Nod
}

func NewNod(id string, message Message) (*Nod, error) {
	if id == "" {
		return nil, fmt.Errorf("nod id cannot be empty")
	}

	return &Nod{
		id:      id,
		message: message,
		next:    nil, // Initialize with no next node
	}, nil
}

func (n *Nod) ID() string {
	return n.id
}

func (n *Nod) Next() *Nod {
	return n.next
}

func (n *Nod) SetNext(next *Nod) {
	n.next = next
}

func (n *Nod) Message() Message {
	return n.message
}

func (n *Nod) SetMessage(message string) {
	n.message = Message(message)
}

func (n *Nod) Bang(stack forth.Stack, state forth.State, modifier string) (forth.Stack, forth.State, []string, error) {
	msg := string(n.message)

	if modifier != "" {
		msg = fmt.Sprintf("`%s` %s", n.message, modifier)
	}

	// This allow for nodTime substitutions
	msgWithSigils, err := forth.ParseSigils(msg)

	if err != nil {
		return stack, state, nil, fmt.Errorf("error parsing sigil")
	}

	newStack, newState, output := forth.Interpret(msgWithSigils, stack, state)

	if len(output) > 0 && strings.HasPrefix(output[0], "Error:") {
		return newStack, newState, nil, fmt.Errorf(output[0])
	}

	return newStack, newState, output, nil
}
