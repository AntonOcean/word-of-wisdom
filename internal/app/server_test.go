package app_test

import (
	"bufio"
	"errors"
	"github.com/stretchr/testify/assert"
	"net"
	"sync"
	"testing"
	"time"
	"word-of-wisdom/internal/app"
	"word-of-wisdom/internal/config"
	"word-of-wisdom/pkg/logger"
)

// MockHandler simulates request handling.
type MockHandler struct{}

func (m *MockHandler) HandleConnection(_ app.Conn) error {
	// Simulate processing delay
	time.Sleep(100 * time.Millisecond)
	return nil
}

// MockHandlerWithError simulates a failing handler
type MockHandlerWithError struct{}

func (m *MockHandlerWithError) HandleConnection(_ app.Conn) error {
	return errors.New("mock handler error")
}

// MockHandler simulates request handling.
type MockHandlerWithPanic struct{}

func (m *MockHandlerWithPanic) HandleConnection(_ app.Conn) error {
	// Simulate processing delay
	panic("hello panic")

}

// TestServerLifecycle tests server start and graceful shutdown.
func TestServerLifecycle(t *testing.T) {
	port := "localhost:8081"

	cfg := config.Config{
		Port:                port,
		MaxConnections:      100,
		ConnectionTimeout:   5 * time.Second,
		ShutdownTimeout:     5 * time.Second,
		RateLimitEvery100MS: 5,
	}

	server := app.NewServer(cfg, logger.GetLogger(), &MockHandler{})

	go server.Start()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Ensure server is running by trying to connect
	conn, err := net.Dial("tcp", port)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	conn.Close()

	// Trigger shutdown
	server.Shutdown()

	// Ensure server stopped accepting connections
	_, err = net.Dial("tcp", port)
	if err == nil {
		t.Fatal("Expected server to be shut down, but it still accepts connections")
	}
}

// TestConnectionHandling checks if the server correctly handles a client request.
func TestConnectionHandling(t *testing.T) {
	port := "localhost:8082"

	cfg := config.Config{
		Port:                port,
		MaxConnections:      100,
		ConnectionTimeout:   5 * time.Second,
		ShutdownTimeout:     5 * time.Second,
		RateLimitEvery100MS: 5,
	}

	server := app.NewServer(cfg, logger.GetLogger(), &MockHandler{})

	go server.Start()
	defer server.Shutdown()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	conn, err := net.Dial("tcp", port)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	conn.Close()
}

// TestConnectionLimit ensures the server enforces maxConnections.
func TestConnectionLimit(t *testing.T) {
	port := "localhost:8083"
	maxConnections := 2

	cfg := config.Config{
		Port:                port,
		MaxConnections:      maxConnections,
		ConnectionTimeout:   5 * time.Second,
		ShutdownTimeout:     5 * time.Second,
		RateLimitEvery100MS: 5,
	}

	server := app.NewServer(cfg, logger.GetLogger(), &MockHandler{})

	go server.Start()
	defer server.Shutdown()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	var conns []net.Conn

	for i := 0; i < maxConnections; i++ {
		conn, err := net.Dial("tcp", port)
		if err != nil {
			t.Fatalf("Failed to connect to server: %v", err)
		}

		conns = append(conns, conn)
	}

	// The last connection should be rejected due to maxConnections limit
	conn, _ := net.Dial("tcp", port)

	// Read response from the server (this is to check if the server rejected the connection)
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err == nil {
		t.Fatal("Expected connection to be rejected due to maxConnections limit")
	}

	conn.Close()

	// Clean up
	for _, c := range conns {
		c.Close()
	}
}

// TestGracefulShutdown checks if the server waits for active connections before shutting down.
func TestGracefulShutdown(t *testing.T) {
	port := "localhost:8084"
	shutdownTimeout := 5 * time.Second

	cfg := config.Config{
		Port:                port,
		MaxConnections:      100,
		ConnectionTimeout:   5 * time.Second,
		ShutdownTimeout:     shutdownTimeout,
		RateLimitEvery100MS: 5,
	}

	server := app.NewServer(cfg, logger.GetLogger(), &MockHandler{})

	go server.Start()
	time.Sleep(100 * time.Millisecond) // Give server time to start

	conn, err := net.Dial("tcp", port)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Shutdown()
	}()

	// Ensure shutdown waits for active connections
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Shutdown completed successfully
	case <-time.After(shutdownTimeout + 500*time.Millisecond):
		t.Fatal("Server did not shut down gracefully within timeout")
	}
}

// TestHandlerError ensures handler errors are logged but do not crash the server
func TestHandlerError(t *testing.T) {
	port := "localhost:8087"

	cfg := config.Config{
		Port:                port,
		MaxConnections:      100,
		ConnectionTimeout:   5 * time.Second,
		ShutdownTimeout:     5 * time.Second,
		RateLimitEvery100MS: 5,
	}

	server := app.NewServer(cfg, logger.GetLogger(), &MockHandlerWithError{})

	go server.Start()
	defer server.Shutdown()

	time.Sleep(100 * time.Millisecond)

	conn, err := net.Dial("tcp", port)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	conn.Close()
}

// TestConnectionRejectionOnShutdown ensures new connections are rejected during shutdown
func TestConnectionRejectionOnShutdown(t *testing.T) {
	port := "localhost:8086"

	cfg := config.Config{
		Port:                port,
		MaxConnections:      100,
		ConnectionTimeout:   5 * time.Second,
		ShutdownTimeout:     5 * time.Second,
		RateLimitEvery100MS: 5,
	}

	server := app.NewServer(cfg, logger.GetLogger(), &MockHandler{})

	go server.Start()
	time.Sleep(100 * time.Millisecond)

	go func() {
		time.Sleep(50 * time.Millisecond)
		server.Shutdown()
	}()

	time.Sleep(100 * time.Millisecond)

	_, err := net.Dial("tcp", port)
	if err == nil {
		t.Fatal("Expected connection to be rejected after shutdown")
	}
}

// TestMultipleClients checks if multiple clients can connect concurrently
func TestMultipleClients(t *testing.T) {
	port := "localhost:8085"

	cfg := config.Config{
		Port:                port,
		MaxConnections:      100,
		ConnectionTimeout:   5 * time.Second,
		ShutdownTimeout:     5 * time.Second,
		RateLimitEvery100MS: 5,
	}

	server := app.NewServer(cfg, logger.GetLogger(), &MockHandler{})

	go server.Start()
	defer server.Shutdown()

	time.Sleep(100 * time.Millisecond) // Give server time to start

	var wg sync.WaitGroup
	clientCount := 5

	for i := 0; i < clientCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			conn, err := net.Dial("tcp", port)
			if err != nil {
				t.Errorf("Client %d failed to connect: %v", i, err)
				return
			}
			conn.Close()
		}(i)
	}

	wg.Wait()
}

// TestPanicRecovery ensures that panics are caught and logged
func TestPanicRecovery(t *testing.T) {
	port := "localhost:8088"

	cfg := config.Config{
		Port:                port,
		MaxConnections:      100,
		ConnectionTimeout:   5 * time.Second,
		ShutdownTimeout:     5 * time.Second,
		RateLimitEvery100MS: 5,
	}

	server := app.NewServer(cfg, logger.GetLogger(), &MockHandlerWithPanic{})

	go server.Start()
	defer server.Shutdown()

	time.Sleep(100 * time.Millisecond) // Give server time to start

	conn, err := net.Dial("tcp", port)
	assert.NoError(t, err, "Client should be able to connect")
	defer conn.Close()

	response, err := bufio.NewReader(conn).ReadString('\n')
	assert.NoError(t, err, "Should receive response from server")
	assert.Equal(t, app.MsgOnErrInternal, response, "Server should handle panics gracefully")
}

// TestRateLimiting ensures that rate limiting works as expected
func TestRateLimiting(t *testing.T) {
	port := "localhost:8089"

	cfg := config.Config{
		Port:                port,
		MaxConnections:      100,
		ConnectionTimeout:   5 * time.Second,
		ShutdownTimeout:     5 * time.Second,
		RateLimitEvery100MS: 2,
	}

	server := app.NewServer(cfg, logger.GetLogger(), &MockHandler{})

	go server.Start()
	defer server.Shutdown()

	time.Sleep(100 * time.Millisecond) // Give server time to start

	// Simulate a client exceeding the rate limit
	conn1, _ := net.Dial("tcp", port)
	conn2, _ := net.Dial("tcp", port)
	conn3, _ := net.Dial("tcp", port)

	// Read responses
	buf1 := bufio.NewReader(conn1)
	buf2 := bufio.NewReader(conn2)
	buf3 := bufio.NewReader(conn3)

	res1, _ := buf1.ReadString('\n')
	res2, _ := buf2.ReadString('\n')
	res3, _ := buf3.ReadString('\n')

	assert.Equal(t, "", res1)
	assert.Equal(t, "", res2)
	assert.Equal(t, app.MsgOnManyReq, res3)

	conn1.Close()
	conn2.Close()
	conn3.Close()
}
