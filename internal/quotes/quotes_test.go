package quotes_test

import (
	"testing"
	"word-of-wisdom/internal/quotes"
)

// TestRandomQuoteProvider ensures GetQuote returns a valid quote from the predefined list.
func TestRandomQuoteProvider(t *testing.T) {
	q := []string{
		"The only limit to our realization of tomorrow is our doubts of today.",
		"Do what you can, with what you have, where you are.",
		"The journey of a thousand miles begins with one step.",
		"Opportunities don't happen. You create them.",
	}

	provider := quotes.NewRandomQuoteProvider(q)
	quotesSet := make(map[string]bool)

	// Populate the set with known quotes
	for i := range q {
		quotesSet[q[i]] = true
	}

	// Check multiple calls return a valid quote
	for i := 0; i < 10; i++ {
		quote := provider.GetQuote()
		if !quotesSet[quote] {
			t.Errorf("Unexpected quote: %s", quote)
		}
	}
}

// TestEmptyQuoteProvider ensures GetQuote doesn't panic when no quotes are available.
func TestEmptyQuoteProvider(t *testing.T) {
	provider := quotes.NewRandomQuoteProvider([]string{})

	quote := provider.GetQuote()
	if quote != quotes.Stub {
		t.Errorf("Expected empty quote, got: %s", quote)
	}
}
