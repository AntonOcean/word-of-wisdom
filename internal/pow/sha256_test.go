package pow_test

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
	"time"
	"word-of-wisdom/internal/pow"
)

// TestGenerateChallenge ensures that the challenge is not empty and varies across calls.
func TestGenerateChallenge(t *testing.T) {
	difficulty := 4

	p := pow.NewSHA256PoW(difficulty)

	challenge1 := p.GenerateChallenge()
	challenge2 := p.GenerateChallenge()

	if challenge1 == "" || challenge2 == "" {
		t.Fatal("Generated challenge should not be empty")
	}
	if challenge1 == challenge2 {
		t.Fatal("Generated challenges should be different")
	}
}

// solvePoW finds a valid solution for a given challenge and difficulty.
func solvePoW(challenge string, difficulty int) string {
	prefix := strings.Repeat("0", difficulty)
	for nonce := 0; ; nonce++ {
		hash := sha256.Sum256([]byte(challenge + fmt.Sprintf("%d", nonce)))
		hashStr := hex.EncodeToString(hash[:])
		if strings.HasPrefix(hashStr, prefix) {
			return fmt.Sprintf("%d", nonce)
		}
	}
}

// TestValidateChallenge checks if the PoW validation correctly accepts or rejects solutions.
func TestValidateChallenge(t *testing.T) {
	difficulty := 4

	p := pow.NewSHA256PoW(difficulty)
	challenge := p.GenerateChallenge()

	// Find a valid solution
	solution := solvePoW(challenge, difficulty)

	// Validate the solution
	if !p.ValidateChallenge(challenge, solution) {
		t.Fatal("Valid PoW solution was rejected")
	}

	// Test with an invalid solution
	invalidSolution := "invalid"
	if p.ValidateChallenge(challenge, invalidSolution) {
		t.Fatal("Invalid PoW solution was accepted")
	}
}

// TestDifficultyLevel ensures that higher difficulty requires more work.
func TestDifficultyLevel(t *testing.T) {
	powLow := pow.NewSHA256PoW(2)
	powHigh := pow.NewSHA256PoW(5)

	challenge := powLow.GenerateChallenge()
	solutionLow := solvePoW(challenge, 2)
	solutionHigh := solvePoW(challenge, 5)

	// Ensure both solutions are valid
	if !powLow.ValidateChallenge(challenge, solutionLow) {
		t.Fatal("Low difficulty PoW solution was rejected")
	}
	if !powHigh.ValidateChallenge(challenge, solutionHigh) {
		t.Fatal("High difficulty PoW solution was rejected")
	}

	// Ensure higher difficulty requires longer solutions
	if len(solutionHigh) < len(solutionLow) {
		t.Fatal("Higher difficulty should result in longer or harder solutions")
	}
}

// TestEmptyChallenge ensures that an empty challenge is rejected.
func TestEmptyChallenge(t *testing.T) {
	p := pow.NewSHA256PoW(4)
	if p.ValidateChallenge("", "solution") {
		t.Fatal("Empty challenge should not be accepted")
	}
}

// TestExtremeSolutionValues ensures extreme inputs do not pass.
func TestExtremeSolutionValues(t *testing.T) {
	p := pow.NewSHA256PoW(4)
	challenge := p.GenerateChallenge()

	// Edge cases: very long string, special characters
	extremeSolutions := []string{
		strings.Repeat("A", 10000),
		"@#$%^&*()_+{}|:<>?",
		"\n\t",
		"",
	}

	for _, solution := range extremeSolutions {
		if p.ValidateChallenge(challenge, solution) {
			t.Fatalf("Extreme solution %q should not be valid", solution)
		}
	}
}

// TestPerformance ensures PoW validation runs within a reasonable time.
func TestPerformance(t *testing.T) {
	p := pow.NewSHA256PoW(4)
	challenge := p.GenerateChallenge()
	solution := solvePoW(challenge, 4)

	start := time.Now()
	if !p.ValidateChallenge(challenge, solution) {
		t.Fatal("Valid PoW solution was rejected")
	}
	elapsed := time.Since(start)

	// Ensure validation runs quickly (less than 10ms)
	if elapsed > 10*time.Millisecond {
		t.Fatalf("PoW validation took too long: %s", elapsed)
	}
}
