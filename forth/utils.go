package forth

import "fmt"

func GetArray(stack Stack) ([]interface{}, Stack, error) {
	var arr []interface{}

	newStack := make(Stack, len(stack))
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

func PopArray(stack Stack) ([]interface{}, Stack, error) {
	if len(stack) == 0 {
		return nil, stack, fmt.Errorf("stack underflow")
	}
	val := stack[len(stack)-1]
	newStack := stack[:len(stack)-1]
	arr, ok := val.([]interface{})
	if !ok {
		return nil, stack, fmt.Errorf("expected array, got %T", val)
	}
	return arr, newStack, nil
}

func PopInt(stack Stack) (int, Stack, error) {
	if len(stack) == 0 {
		return 0, stack, fmt.Errorf("stack underflow")
	}

	val := stack[len(stack)-1]
	newStack := stack[:len(stack)-1]

	switch v := val.(type) {
	case float64:
		return int(v), newStack, nil
	case int:
		return v, newStack, nil
	default:
		return 0, stack, fmt.Errorf("value is not a number: %v", val)
	}
}

func PopString(stack Stack) (string, Stack, error) {
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
func PopFloat(stack Stack) (float64, Stack, error) {
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
