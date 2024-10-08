package grep

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// matchLine checks if any part of the line matches the pattern
func MatchLine(line []byte, pattern string) (bool, error) {
	if utf8.RuneCountInString(pattern) == 0 {
		return false, fmt.Errorf("unsupported pattern: %q", pattern)
	}

	// Check if pattern starts with '^' and ends with '$'
	if pattern[0] == '^' && pattern[len(pattern)-1] == '$' {
		return matchExact(line, pattern[1:len(pattern)-1]), nil
	}
	// Check if pattern starts with '^'
	if pattern[0] == '^' {
		return matchFromPosition(line, pattern[1:]), nil
	}
	// Check if pattern ends with '$'
	if pattern[len(pattern)-1] == '$' {
		return matchFromEnd(line, pattern[0:len(pattern)-1]), nil
	}

	// Iterate through each position in the line to try matching the pattern
	for i := 0; i < len(line); i++ {
		if strings.HasPrefix(pattern, "[^") && !matchFromPosition(line[i:], pattern) {
			return false, nil
		}
		if matchFromPosition(line[i:], pattern) {
			return true, nil
		}
	}

	return false, nil
}

// matchExact checks if the entire line matches the pattern exactly (for ^...$ patterns)
func matchExact(line []byte, pattern string) bool {
	return matchFromPosition(line, pattern)
}

// matchFromPosition checks if the line from the position specified matches the pattern from given position
func matchFromPosition(line []byte, pattern string) bool {
	i, j := 0, 0
	for i < len(line) && j < len(pattern) {
		lineRune, lineSize := utf8.DecodeRune(line[i:])
		patternRune, patternSize := utf8.DecodeRuneInString(pattern[j:])
		// Handle escape sequences (\d, \w, etc.)
		if patternRune == '\\' {
			j += patternSize
			patternRune, patternSize = utf8.DecodeRuneInString(pattern[j:])

			switch patternRune {
			// Match a digit
			case 'd':
				if !unicode.IsDigit(lineRune) {
					return false
				}
				i += lineSize
				j += patternSize
				continue
			// Match a word character (alphanumeric or underscore)
			case 'w':
				if !unicode.IsLetter(lineRune) && !unicode.IsDigit(lineRune) && lineRune != '_' {
					return false
				}
				i += lineSize
				j += patternSize
				continue
			// Invalid escape sequence
			default:
				return false
			}
		} else if patternRune == '(' {
			// Handle alternation (e.g., (cat|dog))
			j += patternSize
			endIdx := j
			depth := 1 // Keep track of nested parentheses

			for endIdx < len(pattern) {
				if pattern[endIdx] == '(' {
					depth++
				} else if pattern[endIdx] == ')' {
					depth--
					if depth == 0 {
						break
					}
				}
				endIdx++
			}

			if endIdx >= len(pattern) || pattern[endIdx] != ')' {
				return false // Unmatched parentheses
			}

			// Extract the alternation inside the parentheses
			alternation := pattern[j:endIdx]
			alternatives := strings.Split(alternation, "|")

			// Try each alternative
			for _, alt := range alternatives {
				if matchFromPosition(line[i:], alt) {
					// If one of the alternatives matches, skip over the entire `(cat|dog)` part
					j = endIdx + 1
					return matchFromPosition(line[i:], pattern[j:])
				}
			}

			// No alternatives matched
			return false

		} else if patternRune == '[' {
			// Handle character classes like [abc] or [^abc]
			j += patternSize
			isNegated := false
			charSet := ""

			// Check if it's a negated character class
			if nextRune, nextSize := utf8.DecodeRuneInString(pattern[j:]); nextRune == '^' {
				isNegated = true
				j += nextSize // Move past the '^'
			}

			// Collect characters in the character class until we hit a closing `]`
			for {
				char, size := utf8.DecodeRuneInString(pattern[j:])
				if char == ']' {
					j += size // Move past the closing bracket
					break
				}
				charSet += string(char)
				j += size
			}

			// Check if the current lineRune matches the character class
			matches := false
			for _, c := range charSet {
				if lineRune == rune(c) {
					matches = true
					break
				}
			}

			// Handle negated class or normal class behavior
			if (isNegated && matches) || (!isNegated && !matches) {
				// If negated and matched, or not negated but didn't match, return false
				return false
			}

			// If character class matched correctly, proceed
			i += lineSize
			continue
		} // Handle quantifiers like '+' (one or more) and '?' (zero or one)
		if j+patternSize < len(pattern) && (pattern[j+patternSize] == '+' || pattern[j+patternSize] == '?') {
			prevRune := patternRune
			quantifier := pattern[j+patternSize]
			j += patternSize + 1 // Move past the quantifier

			// Handle the '+' quantifier (one or more)
			if quantifier == '+' {
				// Ensure the first occurrence matches
				if prevRune != lineRune {
					return false
				}
				i += lineSize

				// Continue matching as long as we find the same character
				for i < len(line) {
					nextRune, nextSize := utf8.DecodeRune(line[i:])
					if nextRune != prevRune {
						break
					}
					i += nextSize
				}
				continue
			}

			if quantifier == '?' {

				// Case 1: Try skipping the preceding character (zero occurrence)
				if matchFromPosition(line[i:], pattern[j:]) {
					return true
				}

				// Case 2: Try matching the preceding character (one occurrence)
				if lineRune == prevRune {
					i += lineSize
					continue
				}

				// If neither case matches, return false
				return false
			}
		} else if patternRune == '.' {
			// Handle '.' (match any single character)
			i += lineSize
			j += patternSize
			continue
		} else if patternRune != lineRune {
			// Literal character match
			return false
		}

		// Advance both the line and pattern pointers
		i += lineSize
		j += patternSize
	}

	// Ensure the entire pattern was matched
	return j == len(pattern)
}

// matchFromEnd tries to match the pattern at the end of the line
func matchFromEnd(line []byte, pattern string) bool {
	if len(line) < len(pattern) {
		return false
	}
	startIdx := len(line) - len(pattern)
	return matchFromPosition(line[startIdx:], pattern)
}
