// world/hedDictionary.go
package world

import (
	"3body/forth"
	"fmt"
)

// DefineHedDictionary creates forth words specifically for head/node operations
func DefineHedDictionary(memory *Memory2D) map[string]forth.DictionaryWord {
	return map[string]forth.DictionaryWord{

		"hed-new": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			if len(stack) < 2 {
				return stack, state, []string{"Error: stack underflow"}
			}

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

			hed, err := NewHed(
				HedID(int(destX), int(destY)),
				nil,
				nil,
				4,
				"",
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

			stack = append(stack, int(destY))
			stack = append(stack, int(destX))
			return stack, state, nil
		},

		"hed-first": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			if len(stack) < 4 {
				return stack, state, []string{"Error: stack underflow"}
			}

			firstX, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			firstY, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			hedX, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			hedY, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			hed, err := memory.GetHed(hedX, hedY)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error fetching hed: %v", err)}
			}

			nod, err := memory.GetNod(firstX, firstY)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error fetching nod: %v", err)}
			}

			hed.first = nod
			hed.current = nod // Maybe this should be taken care of by the hed structure?

			stack = append(stack, int(hedY))
			stack = append(stack, int(hedX))
			return stack, state, nil
		},

		"hed-last": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			if len(stack) < 4 {
				return stack, state, []string{"Error: stack underflow"}
			}

			lastX, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			lastY, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			hedX, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			hedY, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			hed, err := memory.GetHed(hedX, hedY)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error fetching hed: %v", err)}
			}

			nod, err := memory.GetNod(lastX, lastY)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error fetching nod: %v", err)}
			}

			hed.last = nod

			stack = append(stack, int(hedY))
			stack = append(stack, int(hedX))
			return stack, state, nil
		},

		"hed-wrap": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			if len(stack) < 3 {
				return stack, state, []string{"Error: stack underflow"}
			}

			wrapper, newStack, err := forth.PopString(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error popping wrapper: %v", err)}
			}
			stack = newStack

			hedX, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error popping hedX: %v", err)}
			}
			stack = newStack

			hedY, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error popping hedY: %v", err)}
			}
			stack = newStack

			hed, err := memory.GetHed(hedX, hedY)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error getting hed: %v", err)}
			}

			hed.modifier = wrapper

			stack = append(stack, int(hedY))
			stack = append(stack, int(hedX))
			return stack, state, nil
		},

		// BELOW ARE LEGACY WORDS SORT THROUGH, RENAME, DISCARD
		// If you use any of them in the next couple of weeks then port them to use the above words

		// ( nodY nodX destY destX every -- nedY hedX ) creates a new hed with default values
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
				nil,
				int(every),
				"",
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
		// nodY nodX destY destX wrapperString every hed-wrapped
		"hed-wrapped": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			if len(stack) < 6 {
				return stack, state, []string{"Error: stack underflow"}
			}

			every, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			wrapper, newStack, err := forth.PopString(stack)
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
				nil,
				int(every),
				wrapper,
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

		// destY destX firstY firstX lastY lastX address every
		"hed-loop": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
			if len(stack) < 8 {
				return stack, state, []string{"Error: stack underflow"}
			}

			every, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			address, newStack, err := forth.PopString(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			lastX, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			lastY, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			firstX, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			firstY, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			hedX, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			hedY, newStack, err := forth.PopInt(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error: %v", err)}
			}
			stack = newStack

			firstNod, err := memory.GetNod(firstX, firstY)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error getting nod: %v", err)}
			}

			lastNod, err := memory.GetNod(lastX, lastY)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error getting nod: %v", err)}
			}

			formattedAddress := fmt.Sprintf(`"%s" m-osc`, address)

			hed, err := NewHed(
				HedID(hedX, hedY),
				firstNod,
				lastNod,
				every,
				formattedAddress,
				state,
			)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error creating hed: %v", err)}
			}

			if err := memory.AddHed(hedX, hedY, hed); err != nil {
				return stack, state, []string{fmt.Sprintf("Error adding head: %v", err)}
			}

			stack = append(stack, hedY)
			stack = append(stack, hedX)
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

		"hed-freq": func(stack forth.Stack, state forth.State) (forth.Stack, forth.State, []string) {
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
	}
}
