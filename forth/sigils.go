package forth

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func processSigil(sigilType string, value string) (string, error) {
	switch sigilType {
	case "r":
		parts := strings.Split(value, ":")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid range format: expected format <number>:<number>, got %q", value)
		}
		start, err := strconv.Atoi(parts[0])
		if err != nil {
			return "", fmt.Errorf("invalid start number: %v", err)
		}
		end, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("invalid end number: %v", err)
		}
		if start > end {
			start, end = end, start
		}
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		return strconv.Itoa(r.Intn((end - start + 1)) + start), nil

	// Easy to add new sigil types:
	// case "d":
	//     sides, err := strconv.Atoi(value)
	//     if err != nil {
	//         return "", fmt.Errorf("invalid dice sides: %v", err)
	//     }
	//     r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//     return strconv.Itoa(r.Intn(sides) + 1), nil

	default:
		return "", fmt.Errorf("unknown sigil type %q", sigilType)
	}
}

func ParseSigils(input string) (string, error) {
	words := strings.Split(input, " ")
	for i, word := range words {
		if strings.HasPrefix(word, "$") {
			sigilType := word[1:2]
			value := word[2:]

			result, err := processSigil(sigilType, value)
			if err != nil {
				return "", fmt.Errorf("error processing sigil %q: %v", word, err)
			}
			words[i] = result
		}
	}
	return strings.Join(words, " "), nil
}
