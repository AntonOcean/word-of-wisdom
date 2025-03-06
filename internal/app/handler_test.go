package app_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sync"
	"testing"
	"word-of-wisdom/internal/app"
	"word-of-wisdom/internal/app/mocks"
)

func TestHandleConnection_ValidPoW(t *testing.T) {
	quote := "The only limit to our realization of tomorrow is our doubts of today."

	// Prepare mocks
	mockQuoteProvider := mocks.NewQuoteProvider(t)
	mockQuoteProvider.EXPECT().
		GetQuote().
		Return(quote)

	mockPoW := mocks.NewPowChallenge(t)
	mockPoW.EXPECT().
		GenerateChallenge().
		Return("challenge-1234")

	mockPoW.EXPECT().
		ValidateChallenge("challenge-1234", "solution-1234").
		Return(true)

	// Create handler with mocks
	handler := app.NewHandler(mockQuoteProvider, mockPoW)

	// Create mock connection
	mockConn := mocks.NewConn(t)

	mockConn.EXPECT().
		Write(mock.Anything).
		Return(0, nil)

	mockConn.On("Read", mock.Anything).Return(func(p []byte) int {
		copy(p, "solution-1234\n")
		return len("solution-1234\n")
	}, nil)

	err := handler.HandleConnection(mockConn)
	assert.NoError(t, err)

	// Verify PoW validation was called
	mockConn.AssertExpectations(t)
	mockPoW.AssertExpectations(t)
	mockQuoteProvider.AssertExpectations(t)
}

func TestHandleConnection_InvalidPoW(t *testing.T) {
	// Prepare mocks
	mockQuoteProvider := mocks.NewQuoteProvider(t)

	mockPoW := mocks.NewPowChallenge(t)
	mockPoW.EXPECT().
		GenerateChallenge().
		Return("challenge-1234")
	mockPoW.EXPECT().
		ValidateChallenge("challenge-1234", "invalid-solution").
		Return(false)

	// Create handler with mocks
	handler := app.NewHandler(mockQuoteProvider, mockPoW)

	// Create mock connection
	mockConn := mocks.NewConn(t)

	mockConn.EXPECT().
		Write(mock.Anything).
		Return(0, nil)

	mockConn.On("Read", mock.Anything).Return(func(p []byte) int {
		copy(p, "invalid-solution\n")
		return len("invalid-solution\n")
	}, nil)

	err := handler.HandleConnection(mockConn)
	assert.NoError(t, err)

	// Verify PoW validation was called
	mockConn.AssertExpectations(t)
	mockPoW.AssertExpectations(t)
	mockQuoteProvider.AssertNotCalled(t, "GetQuote")
}

func TestHandleConnection_SendMessageError(t *testing.T) {
	// Prepare mocks
	mockQuoteProvider := mocks.NewQuoteProvider(t)

	mockPoW := mocks.NewPowChallenge(t)
	mockPoW.EXPECT().
		GenerateChallenge().
		Return("challenge-1234")

	// Create handler with mocks
	handler := app.NewHandler(mockQuoteProvider, mockPoW)

	// Create mock connection that returns an error on Write
	mockConn := mocks.NewConn(t)

	mockConn.EXPECT().
		Write(mock.Anything).
		Return(0, fmt.Errorf("write error"))

	// Test send message error
	err := handler.HandleConnection(mockConn)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send message")
}

// Test empty client response (edge case)
func TestHandleConnection_EmptyResponse(t *testing.T) {
	mockQuoteProvider := mocks.NewQuoteProvider(t)

	mockPoW := mocks.NewPowChallenge(t)
	mockPoW.EXPECT().
		GenerateChallenge().
		Return("challenge-1234")
	mockPoW.EXPECT().
		ValidateChallenge("challenge-1234", "").
		Return(false)

	handler := app.NewHandler(mockQuoteProvider, mockPoW)

	mockConn := mocks.NewConn(t)

	mockConn.EXPECT().
		Write(mock.Anything).
		Return(0, nil)

	mockConn.On("Read", mock.Anything).Return(func(p []byte) int {
		copy(p, "\n")
		return len("\n")
	}, nil)

	err := handler.HandleConnection(mockConn)
	assert.NoError(t, err)

	mockConn.AssertExpectations(t)
	mockPoW.AssertExpectations(t)
	mockQuoteProvider.AssertNotCalled(t, "GetQuote")
}

// Test network read failure
func TestHandleConnection_ReadError(t *testing.T) {
	mockQuoteProvider := mocks.NewQuoteProvider(t)

	mockPoW := mocks.NewPowChallenge(t)
	mockPoW.EXPECT().
		GenerateChallenge().
		Return("challenge-1234")

	handler := app.NewHandler(mockQuoteProvider, mockPoW)

	mockConn := mocks.NewConn(t)

	mockConn.EXPECT().
		Write(mock.Anything).
		Return(0, nil)

	mockConn.EXPECT().
		Read(mock.Anything).
		Return(0, fmt.Errorf("read error"))

	err := handler.HandleConnection(mockConn)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read client response")

	mockPoW.AssertExpectations(t)
	mockQuoteProvider.AssertNotCalled(t, "GetQuote")
}

// Test concurrent clients
func TestHandleConnection_ConcurrentClients(t *testing.T) {
	quote := "The only limit to our realization of tomorrow is our doubts of today."

	// Prepare mocks
	mockQuoteProvider := mocks.NewQuoteProvider(t)
	mockQuoteProvider.EXPECT().
		GetQuote().
		Return(quote)

	mockPoW := mocks.NewPowChallenge(t)
	mockPoW.EXPECT().
		GenerateChallenge().
		Return("challenge-1234")

	mockPoW.EXPECT().
		ValidateChallenge("challenge-1234", "solution-1234").
		Return(true)

	handler := app.NewHandler(mockQuoteProvider, mockPoW)

	const numClients = 5
	var wg sync.WaitGroup
	wg.Add(numClients)

	for i := 0; i < numClients; i++ {
		go func() {
			defer wg.Done()

			mockConn := mocks.NewConn(t)
			mockConn.EXPECT().
				Write(mock.Anything).
				Return(0, nil)

			mockConn.On("Read", mock.Anything).Return(func(p []byte) int {
				copy(p, "solution-1234\n")
				return len("solution-1234\n")
			}, nil)

			err := handler.HandleConnection(mockConn)
			assert.NoError(t, err)
		}()
	}

	wg.Wait()
	mockPoW.AssertExpectations(t)
	mockQuoteProvider.AssertExpectations(t)
}
