// world/worldDictionary.go
package world

import (
	"fmt"
	"math/rand"

	"lofl/forth"

	"github.com/hypebeast/go-osc/osc"
)

// DefineWorldDictionary creates forth words that interact with the world
func DefineWorldDictionary(memory *Memory2D, clock *Clock, client *osc.Client, messageChan chan string) map[string]forth.DictionaryWord {
	return map[string]forth.DictionaryWord{
		// random ( -- n ) places a random number on stack
		"random": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			return append(stack, rand.Float64()), state, nil
		},

		// print-memory ( -- ) prints the memory state
		"print-memory": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			fmt.Printf("%+v\n", memory)
			return stack, state, nil
		},

		// start-clock ( -- ) starts the clock
		"start-clock": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			if err := clock.Start(memory); err != nil {
				return stack, state, []string{fmt.Sprintf("Error starting clock: %v", err)}
			}
			return stack, state, nil
		},

		// stop-clock ( -- ) stops the clock
		"stop-clock": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			if err := clock.Stop(); err != nil {
				return stack, state, []string{fmt.Sprintf("Error stopping clock: %v", err)}
			}
			return stack, state, nil
		},

		// m ( message address -- ) sends a message
		"m": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			if len(stack) < 2 {
				return stack, state, []string{"Error: stack underflow"}
			}

			address, newStack, err := popString(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			message, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			oscMsg := osc.NewMessage(fmt.Sprintf("/%s", address))
			oscMsg.Append(float32(message))
			client.Send(oscMsg)

			return stack, state, nil
		},

		// seq ( messages... y x length -- y x ) creates a sequence
		"seq": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			if len(stack) < 3 {
				return stack, state, []string{"Error: stack underflow"}
			}

			length, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			x, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			y, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			messages := make([]string, 0, int(length))
			for i := 0; i < int(length); i++ {
				msg, newStack, err := popString(stack)
				if err != nil {
					return stack, state, []string{fmt.Sprintf("Error getting message: %v", err)}
				}
				stack = newStack
				messages = append([]string{msg}, messages...) // Prepend to reverse order
			}

			err = createSequence(memory, int(x), int(y), messages)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error creating sequence: %v", err)}
			}

			// Push coordinates back on stack
			stack = append(stack, float64(y))
			stack = append(stack, float64(x))

			return stack, state, nil
		},

		// qsm ( array address every y x -- y x ) creates a sequence of messages with a head
		"qsm": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			if len(stack) < 5 {
				return stack, state, []string{"Error: stack underflow"}
			}

			x, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			y, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			every, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			address, newStack, err := popString(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			arr, newStack, err := getArray(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error getting array: %v", err)}
			}
			stack = newStack

			messages := make([]string, len(arr))
			for i, msg := range arr {
				switch v := msg.(type) {
				case string:
					if v == "_" {
						messages[i] = "_"
					} else {
						messages[i] = fmt.Sprintf(`%s "%s" m`, v, address)
					}
				case float64:
					messages[i] = fmt.Sprintf(`%v "%s" m`, v, address)
				case int:
					messages[i] = fmt.Sprintf(`%d "%s" m`, v, address)
				default:
					return stack, state, []string{fmt.Sprintf("Error: array contains invalid type: %T", msg)}
				}
			}

			for i := 0; i < len(messages)/2; i++ {
				j := len(messages) - 1 - i
				messages[i], messages[j] = messages[j], messages[i]
			}

			hed, err := NewHed(
				HedID(int(x), int(y)),
				nil,
				int(every),
				state,
			)

			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error creating head: %v", err)}
			}

			if err := memory.AddHed(int(x), int(y), hed); err != nil {
				return stack, state, []string{fmt.Sprintf("Error adding head: %v", err)}
			}

			if err := createSequence(memory, int(x)+1, int(y), messages); err != nil {
				return stack, state, []string{fmt.Sprintf("Error creating sequence: %v", err)}
			}

			firstNod, err := memory.GetNod(int(x)+1, int(y))
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error getting first node: %v", err)}
			}
			hed.first = firstNod
			hed.current = firstNod

			stack = append(stack, float64(y))
			stack = append(stack, float64(x))
			return stack, state, nil
		},
		// maybe ( message probability -- ) executes message with probability
		"maybe": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			if len(stack) < 2 {
				return stack, state, []string{"Error: stack underflow"}
			}

			prob, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			msg, newStack, err := popString(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			if rand.Float64() < prob {
				return forth.Interpret(msg, stack, state)
			}

			return stack, state, nil
		},

		// one-of ( message2 message1 probability -- ) executes one of two messages
		"one-of": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			if len(stack) < 3 {
				return stack, state, []string{"Error: stack underflow"}
			}

			prob, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			msg1, newStack, err := popString(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			msg2, newStack, err := popString(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			if rand.Float64() < prob {
				return forth.Interpret(msg1, stack, state)
			}
			return forth.Interpret(msg2, stack, state)
		},

		"_": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			return stack, state, nil
		},

		// mod ( y x modMessage -- y x ) adds a modifier to a head
		"mod": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			if len(stack) < 3 {
				return stack, state, []string{"Error: stack underflow"}
			}

			modMsg, newStack, err := popString(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			x, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			y, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			hed, err := memory.GetHed(int(x), int(y))
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error getting head: %v", err)}
			}

			if modMsg == "0" {
				hed.SetModifier("")
			} else {
				hed.SetModifier(modMsg)
			}

			stack = append(stack, float64(y))
			stack = append(stack, float64(x))

			return stack, state, nil
		},
		"qs": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			if len(stack) < 3 {
				return stack, state, []string{"Error: stack underflow"}
			}

			x, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			y, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			every, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			arr, newStack, err := getArray(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error getting array: %v", err)}
			}
			stack = newStack

			messages := make([]string, len(arr))
			for i, msg := range arr {
				messages[i] = fmt.Sprint(msg)
			}

			// Reverse messages
			for i := 0; i < len(messages)/2; i++ {
				j := len(messages) - 1 - i
				messages[i], messages[j] = messages[j], messages[i]
			}

			hed, err := NewHed(
				HedID(int(x), int(y)),
				nil,
				int(every),
				state,
			)

			if err := memory.AddHed(int(x), int(y), hed); err != nil {
				return stack, state, []string{fmt.Sprintf("Error adding head: %v", err)}
			}

			if err := createSequence(memory, int(x)+1, int(y), messages); err != nil {
				return stack, state, []string{fmt.Sprintf("Error creating sequence: %v", err)}
			}

			firstNod, err := memory.GetNod(int(x)+1, int(y))
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error getting first node: %v", err)}
			}

			hed.first = firstNod
			hed.current = firstNod

			stack = append(stack, float64(y))
			stack = append(stack, float64(x))
			return stack, state, nil
		},

		"hed": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			if len(stack) < 5 {
				return stack, state, []string{"Error: stack underflow"}
			}

			every, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			nodX, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			nodY, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			destX, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			destY, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			nod, err := memory.GetNod(int(nodX), int(nodY))
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error getting nod: %v", err)}
			}

			hed, err := NewHed(
				HedID(int(destX), int(destY)),
				nod,
				int(every),
				state,
			)

			if err := memory.AddHed(int(destX), int(destY), hed); err != nil {
				return stack, state, []string{fmt.Sprintf("Error adding head: %v", err)}
			}

			if err := memory.AddHed(int(destX), int(destY), hed); err != nil {
				return stack, state, []string{fmt.Sprintf("Error adding head: %v", err)}
			}

			stack = append(stack, float64(destY))
			stack = append(stack, float64(destX))
			return stack, state, nil
		},

		"nod": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			if len(stack) < 4 {
				return stack, state, []string{"Error: stack underflow"}
			}

			nextX, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			nextY, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			nodX, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			nodY, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			destNod, err := memory.GetNod(int(nextX), int(nextY))
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error fetching destNod: %v", err)}
			}

			nod, err := NewNod(
				NodID(int(nodX), int(nodY)),
				"print",
			)

			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}

			nod.SetNext(destNod)

			if err := memory.AddNod(int(nodX), int(nodY), nod); err != nil {
				return stack, state, []string{fmt.Sprintf("Error adding nod: %v", err)}
			}

			stack = append(stack, float64(nextY))
			stack = append(stack, float64(nextX))
			return stack, state, nil
		},

		"start": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			if len(stack) < 2 {
				return stack, state, []string{"Error: stack underflow"}
			}

			x, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			y, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			hed, err := memory.GetHed(int(x), int(y))
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error getting head: %v", err)}
			}

			hed.Start()

			stack = append(stack, float64(y))
			stack = append(stack, float64(x))
			return stack, state, nil
		},

		"stop": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			if len(stack) < 2 {
				return stack, state, []string{"Error: stack underflow"}
			}

			x, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			y, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			hed, err := memory.GetHed(int(x), int(y))
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error getting head: %v", err)}
			}

			hed.Stop()

			stack = append(stack, float64(y))
			stack = append(stack, float64(x))
			return stack, state, nil
		},

		"point": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			if len(stack) < 4 {
				return stack, state, []string{"Error: stack underflow"}
			}

			x2, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			y2, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			x1, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			y1, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			nod, err := memory.GetNod(int(x1), int(y1))
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error getting nod: %v", err)}
			}

			nextNod, err := memory.GetNod(int(x2), int(y2))
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error getting nod: %v", err)}
			}

			if x1 == x2 && y1 == y2 {
				nod.SetNext(nil)
			} else {
				nod.SetNext(nextNod)
			}

			stack = append(stack, float64(y2))
			stack = append(stack, float64(x2))
			return stack, state, nil
		},

		"r-m": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			if len(stack) < 3 {
				return stack, state, []string{"Error: stack underflow"}
			}

			message, newStack, err := popString(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			x, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			y, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			nod, err := memory.GetNod(int(x), int(y))
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error getting nod: %v", err)}
			}

			nod.SetMessage(message)

			stack = append(stack, float64(y))
			stack = append(stack, float64(x))
			return stack, state, nil
		},

		"r-f": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			if len(stack) < 3 {
				return stack, state, []string{"Error: stack underflow"}
			}

			freq, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			x, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			y, newStack, err := popNumber(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			hed, err := memory.GetHed(int(x), int(y))
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error getting head: %v", err)}
			}

			hed.SetEvery(int(freq))

			stack = append(stack, float64(y))
			stack = append(stack, float64(x))
			return stack, state, nil
		},
		"mg": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			if len(stack) < 1 {
				return stack, state, []string{"Error: stack underflow"}
			}

			msg, newStack, err := popString(stack)

			fmt.Print(msg)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}

			if messageChan != nil {
				messageChan <- msg
			}

			return newStack, state, nil
		},
	}

}
