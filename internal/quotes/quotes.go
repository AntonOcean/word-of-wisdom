package quotes

//go:generate ifacemaker -f quotes.go -s RandomQuoteProvider -p quotes -i QuoteProvider -o interface_generated.go

import (
	"math/rand"
	"time"
)

const Stub = "Angry people are not always wise."

type RandomQuoteProvider struct {
	quotes []string
	rng    *rand.Rand
}

func NewRandomQuoteProvider(quotes []string) QuoteProvider {
	return &RandomQuoteProvider{
		quotes: quotes,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GetQuote returns a random quote from the predefined list
func (q *RandomQuoteProvider) GetQuote() string {
	if len(q.quotes) == 0 {
		return Stub
	}

	return q.quotes[q.rng.Intn(len(q.quotes))]
}
