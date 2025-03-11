package app

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"word-of-wisdom/pkg/protocol"
)

const InvalidMsg = "Invalid PoW solution"

type H struct {
	quoteProvider quoteProvider
	powChallenge  powChallenge
}

func NewHandler(quoteProvider quoteProvider, powChallenge powChallenge) Handler {
	return &H{
		quoteProvider: quoteProvider,
		powChallenge:  powChallenge,
	}
}

// sendMessage sends a message to the client and logs errors.
func sendMessage(conn Conn, message string) error {
	_, err := fmt.Fprintln(conn, message)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// HandleConnection manages a single client connection and performs PoW validation.
func (h *H) HandleConnection(conn Conn) error {
	// Generate and send PoW challenge
	challenge := h.powChallenge.GenerateChallenge()
	if err := sendMessage(conn, protocol.PrefixChallenge+challenge); err != nil {
		return fmt.Errorf("failed to send challenge: %w", err)
	}

	// Read and validate client response
	solution, err := readClientResponse(conn)
	if err != nil {
		return fmt.Errorf("failed to read client response: %w", err)
	}

	// Validate Proof of Work (PoW)
	if !h.powChallenge.ValidateChallenge(challenge, solution) {
		if err := sendMessage(conn, protocol.PrefixError+InvalidMsg); err != nil {
			return fmt.Errorf("failed to send validate: %w", err)
		}

		return nil
	}

	// Send quote if PoW is valid
	quote := h.quoteProvider.GetQuote()
	if err := sendMessage(conn, protocol.PrefixQuote+quote); err != nil {
		return fmt.Errorf("failed to send quote: %w", err)
	}

	return nil
}

// readClientResponse reads the client’s PoW solution from the connection
func readClientResponse(conn Conn) (string, error) {
	const maxReadSize = 1024

	limitedReader := io.LimitedReader{R: conn, N: maxReadSize}

	reader := bufio.NewReader(&limitedReader)
	solution, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(solution), nil
}
