package forth

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// CreateStack initializes an empty stack
func CreateStack() Stack {
	return make(Stack, 0)
}

// CreateInitialState initializes interpreter state
func CreateInitialState() State {
	return State{
		Dictionary:        createInitialDictionary(),
		Compiling:         false,
		CurrentDefinition: make([]string, 0),
		CurrentWord:       nil,
	}
}

// Push adds an item to the stack
func Push(stack Stack, item StackItem) Stack {
	return append(stack, item)
}

// Pop removes and returns the top item from the stack
func Pop(stack Stack) (Stack, StackItem, error) {
	if len(stack) == 0 {
		return stack, nil, errors.New("stack underflow")
	}
	item := stack[len(stack)-1]
	return stack[:len(stack)-1], item, nil
}

// createInitialDictionary creates the basic Forth dictionary
func createInitialDictionary() Dictionary {
	return Dictionary{
		"+": func(stack Stack, state State) (Stack, State, []string) {
			if len(stack) < 2 {
				return stack, state, []string{"stack underflow"}
			}
			s, b, _ := Pop(stack)
			s1, a, _ := Pop(s)
			return Push(s1, a.(float64)+b.(float64)), state, nil
		},
		"-": func(stack Stack, state State) (Stack, State, []string) {
			if len(stack) < 2 {
				return stack, state, []string{"stack underflow"}
			}
			s, b, _ := Pop(stack)
			s1, a, _ := Pop(s)
			return Push(s1, a.(float64)-b.(float64)), state, nil
		},
		// Add other basic operations here
		":": func(stack Stack, state State) (Stack, State, []string) {
			if state.Compiling {
				return stack, state, []string{"nested definitions not allowed"}
			}
			newState := state
			newState.Compiling = true
			return stack, newState, nil
		},
		";": func(stack Stack, state State) (Stack, State, []string) {
			if !state.Compiling {
				return stack, state, []string{"not in compilation mode"}
			}
			if state.CurrentWord == nil {
				return stack, state, []string{"no word name provided"}
			}

			// Create new word from current definition
			wordName := *state.CurrentWord
			definition := state.CurrentDefinition

			state.Dictionary[wordName] = func(s Stack, st State) (Stack, State, []string) {
				return Interpret(strings.Join(definition, " "), s, st)
			}

			newState := state
			newState.Compiling = false
			newState.CurrentDefinition = nil
			newState.CurrentWord = nil

			return stack, newState, nil
		},
		"[": func(stack Stack, state State) (Stack, State, []string) {

			return Push(stack, "["), state, nil
		},
		"]": func(stack Stack, state State) (Stack, State, []string) {

			arr, newStack, err := GetArray(stack)
			if err != nil {
				return stack, state, []string{fmt.Sprintf("Error creating array: %v", err)}
			}

			return Push(newStack, arr), state, nil
		},
		"print-array": func(stack Stack, state State) (Stack, State, []string) {
			if len(stack) < 1 {
				return stack, state, []string{"stack underflow"}
			}

			s, item, _ := Pop(stack)

			// Try to convert item to array
			arr, ok := item.([]interface{})
			if !ok {
				return stack, state, []string{"top item is not an array"}
			}

			// Build string representation of array
			var elements []string
			for _, val := range arr {
				switch v := val.(type) {
				case float64:
					elements = append(elements, fmt.Sprintf("%g", v))
				case string:
					elements = append(elements, fmt.Sprintf("\"%s\"", v))
				default:
					elements = append(elements, fmt.Sprintf("%v", v))
				}
			}

			return s, state, []string{fmt.Sprintf("[ %s ]", strings.Join(elements, " "))}
		},
	}
}

// splitPreservingStrings splits input while preserving quoted strings
func splitPreservingStrings(input string) []string {
	var tokens []string
	var currentToken strings.Builder
	inBacktick := false
	inQuote := false

	for _, char := range input {
		switch {
		case char == '`' && !inQuote:
			if inBacktick {
				currentToken.WriteRune(char)
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
				inBacktick = false
			} else {
				if currentToken.Len() > 0 {
					tokens = append(tokens, currentToken.String())
					currentToken.Reset()
				}
				currentToken.WriteRune(char)
				inBacktick = true
			}
		case char == '"' && !inBacktick:
			if inQuote {
				currentToken.WriteRune(char)
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
				inQuote = false
			} else {
				if currentToken.Len() > 0 {
					tokens = append(tokens, currentToken.String())
					currentToken.Reset()
				}
				currentToken.WriteRune(char)
				inQuote = true
			}
		case char == ' ' && !inBacktick && !inQuote:
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
		default:
			currentToken.WriteRune(char)
		}
	}

	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}

	return tokens
}

// Interpret processes a Forth string and returns the new stack and state
func Interpret(input string, stack Stack, state State) (Stack, State, []string) {
	words := splitPreservingStrings(input)
	currentStack := stack
	currentState := state
	var output []string

	for _, word := range words {
		if currentState.Compiling {
			if currentState.CurrentWord == nil {
				wordName := word
				currentState.CurrentWord = &wordName
				continue
			}
			if word != ";" {
				currentState.CurrentDefinition = append(currentState.CurrentDefinition, word)
				continue
			}
		}

		// sigilnotation
		// these are replaced in place when the nod is evaluated
		if strings.HasPrefix(word, "$") {
			currentStack = Push(currentStack, word)
		} else if dictWord, exists := currentState.Dictionary[word]; exists {
			var newOutput []string
			currentStack, currentState, newOutput = dictWord(currentStack, currentState)
			output = append(output, newOutput...)
		} else if num, err := strconv.ParseFloat(word, 64); err == nil {
			currentStack = Push(currentStack, num)
		} else if (strings.HasPrefix(word, "\"") && strings.HasSuffix(word, "\"")) ||
			(strings.HasPrefix(word, "`") && strings.HasSuffix(word, "`")) {
			currentStack = Push(currentStack, word[1:len(word)-1])
		} else if strings.HasPrefix(word, "'") {
			currentStack = Push(currentStack, word[1:])
		} else {
			return currentStack, currentState, append(output, "Unknown word: "+word)
		}
	}

	return currentStack, currentState, output
}
