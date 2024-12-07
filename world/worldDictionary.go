// world/worldDictionary.go
package world

import (
	"fmt"
	"math/rand"

	"lofl/forth"

	"lofl/connections"

	"github.com/hypebeast/go-osc/osc"
)

// DefineWorldDictionary creates forth words that interact with the world
func DefineWorldDictionary(memory *Memory2D, clock *Clock, client *osc.Client) map[string]forth.DictionaryWord {
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

			address, newStack, err := forth.PopString(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			message, newStack, err := forth.PopFloat(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			oscMsg := osc.NewMessage(fmt.Sprintf("/%s", address))
			oscMsg.Append(float32(message))
			client.Send(oscMsg)

			return stack, state, nil
		},

		// array wrapper y x seq -> sequence of nods / y x of first nod
		"seq": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			x, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			y, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			wrapper, newStack, err := forth.PopString(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			arr, newStack, err := forth.PopArray(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			// Reverses the array and then wraps the strings
			nodes := make([]*Nod, len(arr))
			for i, val := range arr {
				var wrappedMessage string
				switch v := val.(type) {
				case string:
					wrappedMessage = fmt.Sprintf("%s %s", v, wrapper)
				case float64:
					wrappedMessage = fmt.Sprintf("%g %s", v, wrapper)
				case int:
					wrappedMessage = fmt.Sprintf("%d %s", v, wrapper)
				default:
					wrappedMessage = fmt.Sprintf("%v %s", v, wrapper)
				}

				nod, err := NewNod(NodID(x+i, y), Message(wrappedMessage))
				if err != nil {
					return stack, state, []string{fmt.Sprintf("error creating node: %v", err)}
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
					return stack, state, []string{fmt.Sprintf("error adding node: %v", err)}
				}
			}

			stack = forth.Push(stack, y)
			stack = forth.Push(stack, x)
			return stack, state, nil
		},

		//array address every y x -- y x
		// Deprectated, use namespaced qs-m
		"qsm": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {

			return forth.Interpret("qs-m", stack, state)

		},

		"qs-m": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			x, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			y, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			every, newStack, err := forth.PopFloat(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			address, newStack, err := forth.PopString(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			arr, newStack, err := forth.PopArray(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error getting array: %v", err)}
			}
			stack = newStack

			stack = forth.Push(stack, arr)
			stack = forth.Push(stack, fmt.Sprintf("%s m", address)) // push the wrapper
			stack = forth.Push(stack, y)
			stack = forth.Push(stack, x+1)

			stack, state, message := forth.Interpret(fmt.Sprintf("seq %d %d %f hed", y, x, every), stack, state)

			return stack, state, message
		},

		"qs": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			x, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			y, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			every, newStack, err := forth.PopFloat(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			arr, newStack, err := forth.PopArray(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error getting array: %v", err)}
			}
			stack = newStack

			stack = forth.Push(stack, arr)
			stack = forth.Push(stack, "") // push the wrapper
			stack = forth.Push(stack, y)
			stack = forth.Push(stack, x+1)

			stack, state, message := forth.Interpret(fmt.Sprintf("seq %d %d %f hed", y, x, every), stack, state)

			return stack, state, message
		},

		// maybe ( message probability -- ) executes message with probability
		"maybe": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			if len(stack) < 2 {
				return stack, state, []string{"Error: stack underflow"}
			}

			prob, newStack, err := forth.PopFloat(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			msg, newStack, err := forth.PopString(stack)
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

			prob, newStack, err := forth.PopFloat(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			msg1, newStack, err := forth.PopString(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			msg2, newStack, err := forth.PopString(stack)
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

			modMsg, newStack, err := forth.PopString(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			x, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			y, newStack, err := forth.PopInt(stack)
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

		"hed": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			if len(stack) < 5 {
				return stack, state, []string{"Error: stack underflow"}
			}

			every, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			destX, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			destY, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			nodX, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			nodY, newStack, err := forth.PopInt(stack)
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
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error creating hed: %v", err)}
			}

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

			nextX, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			nextY, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			nodX, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			nodY, newStack, err := forth.PopInt(stack)
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

			x, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			y, newStack, err := forth.PopInt(stack)
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

			x, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			y, newStack, err := forth.PopInt(stack)
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

			x2, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			y2, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			x1, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			y1, newStack, err := forth.PopInt(stack)
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

			message, newStack, err := forth.PopString(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			x, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			y, newStack, err := forth.PopInt(stack)
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

			freq, newStack, err := forth.PopFloat(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			x, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			y, newStack, err := forth.PopInt(stack)
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
		"mlg": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			if len(stack) < 1 {
				return stack, state, []string{"Error: stack underflow"}
			}

			msg, newStack, err := forth.PopString(stack)

			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}

			if connections.HTTPMessageChannel != nil {
				connections.HTTPMessageChannel <- connections.HTTPMessage{Type: "line", Content: msg}
			}

			return newStack, state, nil
		},
		"mhg": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			if len(stack) < 1 {
				return stack, state, []string{"Error: stack underflow"}
			}

			msg, newStack, err := forth.PopString(stack)

			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}

			if connections.HTTPMessageChannel != nil {
				connections.HTTPMessageChannel <- connections.HTTPMessage{Type: "hydra", Content: msg}
			}

			return newStack, state, nil
		},
	}

}
