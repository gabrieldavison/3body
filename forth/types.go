package forth

// StackItem represents items that can be stored on the stack
type StackItem interface{}

// Stack is a slice of items
type Stack []StackItem

// DictionaryWord is a function that manipulates the stack and state
type DictionaryWord func(stack Stack, state State) (Stack, State, []string)

// Dictionary maps word names to their implementations
type Dictionary map[string]DictionaryWord

// State maintains the interpreter's state
type State struct {
	Dictionary        Dictionary
	Compiling         bool
	CurrentDefinition []string
	CurrentWord       *string
}
