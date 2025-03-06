package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strings"
	"word-of-wisdom/pkg/protocol"
)

const difficulty = 4 // Match server difficulty

func solvePoW(challenge string) string {
	var solution int64 = 0
	for {
		hash := sha256.Sum256([]byte(fmt.Sprintf("%s%d", challenge, solution)))
		hashStr := hex.EncodeToString(hash[:])
		if strings.HasPrefix(hashStr, strings.Repeat("0", difficulty)) {
			return fmt.Sprintf("%d", solution)
		}
		solution++
	}
}

func main() {
	conn, err := net.Dial("tcp", "wisdom-server:9000") // Server hostname in Docker
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	message, _ := reader.ReadString('\n')
	fmt.Println("Server Message:", message)

	if strings.HasPrefix(message, protocol.PrefixChallenge) {
		challenge := strings.TrimPrefix(message, protocol.PrefixChallenge)
		challenge = strings.TrimSpace(challenge)

		// Solve PoW
		solution := solvePoW(challenge)
		fmt.Fprintf(conn, "%s\n", solution)

		// Read response
		quote, _ := reader.ReadString('\n')
		fmt.Println("Server Response:", quote)
	} else {
		fmt.Println("Unexpected response from server:", message)
	}
}
