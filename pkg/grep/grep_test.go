package grep

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchLine(t *testing.T) {
	var tests = []struct {
		name      string
		line      []byte
		pattern   string
		isMatched bool
	}{
		{
			name:      "Match literal character",
			line:      []byte("apple"),
			pattern:   "a",
			isMatched: true,
		},
		{
			name:      "No literal character match",
			line:      []byte("dog"),
			pattern:   "a",
			isMatched: false,
		},
		{
			name:      "Match digits",
			line:      []byte("apple123"),
			pattern:   "\\d\\d\\d",
			isMatched: true,
		},
		{
			name:      "No digits match",
			line:      []byte("apple"),
			pattern:   "\\d\\d\\d",
			isMatched: false,
		},
		{
			name:      "Match alphanumeric characters",
			line:      []byte("alphanum3ric"),
			pattern:   "\\w",
			isMatched: true,
		},
		{
			name:      "No alphanumeric characters match",
			line:      []byte("..?"),
			pattern:   "\\w",
			isMatched: false,
		},
		{
			name:      "Match positive character groups",
			line:      []byte("apple"),
			pattern:   "[abc]",
			isMatched: true,
		},
		{
			name:      "No positive character groups match",
			line:      []byte("apple"),
			pattern:   "[ghi]",
			isMatched: false,
		},
		{
			name:      "Match negative character groups",
			line:      []byte("dog"),
			pattern:   "[^abc]",
			isMatched: true,
		},
		{
			name:      "No negative character groups match",
			line:      []byte("apple"),
			pattern:   "[^abc]",
			isMatched: false,
		},
		{
			name:      "Match combining character classes",
			line:      []byte("caaaaaats"),
			pattern:   "^ca+ts$",
			isMatched: true,
		},
		{
			name:      "No combining character classes match",
			line:      []byte("1 dog"),
			pattern:   "\\d \\w\\w\\ws",
			isMatched: false,
		},
		{
			name:      "Match start symbol",
			line:      []byte("cats"),
			pattern:   "^cats",
			isMatched: true,
		},
		{
			name:      "No start symbol match",
			line:      []byte("dogs"),
			pattern:   "^cats",
			isMatched: false,
		},
		{
			name:      "Match end symbol",
			line:      []byte("cats"),
			pattern:   "cats$",
			isMatched: true,
		},
		{
			name:      "No end symbol match",
			line:      []byte("cat"),
			pattern:   "cats$",
			isMatched: false,
		},
		{
			name:      "Match + symbol",
			line:      []byte("caaaaaaats"),
			pattern:   "ca+ts",
			isMatched: true,
		},
		{
			name:      "No + symbol match",
			line:      []byte("ct"),
			pattern:   "ca+ts",
			isMatched: false,
		},
		{
			name:      "Match 0 time to question symbol",
			line:      []byte("dog"),
			pattern:   "dog?",
			isMatched: true,
		},
		{
			name:      "Match to question symbol",
			line:      []byte("dogs"),
			pattern:   "dog?",
			isMatched: true,
		},
		{
			name:      "No match to question symbol",
			line:      []byte("c"),
			pattern:   "ca?",
			isMatched: false,
		},
		{
			name:      "Match to . symbol",
			line:      []byte("dogs"),
			pattern:   "dog.",
			isMatched: true,
		},
		{
			name:      "No match to . symbol",
			line:      []byte("c"),
			pattern:   "c.",
			isMatched: false,
		},
		{
			name:      "Match to | symbol",
			line:      []byte("dog|cat"),
			pattern:   "dog",
			isMatched: true,
		},
		{
			name:      "No match to | symbol",
			line:      []byte("cat|dog"),
			pattern:   "bird",
			isMatched: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			matched, _ := MatchLine(test.line, test.pattern)
			assert.Equal(t, test.isMatched, matched)
		})
	}
}
