package pow

//go:generate ifacemaker -f sha256.go -s SHA256PoW -p pow -i PoW -o interface_generated.go

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type SHA256PoW struct {
	difficulty int
	rng        *rand.Rand
}

func NewSHA256PoW(difficulty int) PoW {
	return &SHA256PoW{
		difficulty: difficulty,
		rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GenerateChallenge creates a random challenge string.
func (p *SHA256PoW) GenerateChallenge() string {
	return fmt.Sprintf("%x", p.rng.Int63())
}

// ValidateChallenge checks if the provided solution meets the required difficulty.
func (p *SHA256PoW) ValidateChallenge(challenge, solution string) bool {
	hash := sha256.Sum256([]byte(challenge + solution))
	hashStr := hex.EncodeToString(hash[:]) // TODO improve it with binary
	return strings.HasPrefix(hashStr, strings.Repeat("0", p.difficulty))
}
