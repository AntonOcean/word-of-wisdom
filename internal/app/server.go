package app

import (
	"context"
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
	"word-of-wisdom/internal/config"
)

// Server struct encapsulates the server's state and behavior
type Server struct {
	listener     net.Listener
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	semaphore    chan struct{}
	shutdownOnce sync.Once
	config       config.Config
	handler      Handler
	logger       *logrus.Logger
}

// NewServer creates and initializes a new server instance
func NewServer(c config.Config, logger *logrus.Logger, handler Handler) *Server {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	return &Server{
		ctx:       ctx,
		cancel:    cancel,
		semaphore: make(chan struct{}, c.MaxConnections),
		handler:   handler,
		config:    c,
		logger:    logger,
	}
}

// Start initializes the listener, starts accepting connections, and waits for shutdown
func (s *Server) Start() {
	var err error

	s.listener, err = net.Listen("tcp", s.config.Port)
	if err != nil {
		s.logger.Errorf("Failed to start server: %v", err)
		return
	}

	s.logger.Infof("Server started on port %s", s.config.Port)

	// Start accepting connections
	go s.acceptConnections()

	// Wait for shutdown signal
	<-s.ctx.Done()

	// Perform graceful shutdown
	s.Shutdown()
}

// acceptConnections listens for incoming connections and limits concurrency
func (s *Server) acceptConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.ctx.Done():
				s.logger.Info("Server is shutting down, stopping connection handling...")
				return
			default:
				if strings.Contains(err.Error(), "use of closed network connection") {
					s.logger.Info("Listener closed, stopping connection handling...")
					return
				}

				s.logger.Errorf("Failed to accept connection: %v", err)
			}
			continue
		}

		// Enforce connection limit
		select {
		case s.semaphore <- struct{}{}: // Acquire a slot
			s.wg.Add(1)
			go s.handleClient(conn)
		default:
			s.logger.Info("Too many connections. Rejecting client.")
			conn.Close()
		}
	}
}

// handleClient processes a single client connection and ensures proper resource release
func (s *Server) handleClient(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()
	defer func() { <-s.semaphore }() // Release the slot after processing

	err := conn.SetDeadline(time.Now().Add(s.config.ConnectionTimeout))
	if err != nil {
		s.logger.Errorf("Failed to set deadline: %v", err)
	}

	err = s.handler.HandleConnection(conn) // User-defined function to process client requests
	if err != nil {
		s.logger.Errorf("Error: %v", err)
	}
}

// Shutdown gracefully shuts down the server, closing the listener and waiting for active connections to finish
func (s *Server) Shutdown() {
	s.shutdownOnce.Do(func() {
		s.logger.Info("Shutting down server...")

		// Close the listener to stop accepting new connections
		if err := s.listener.Close(); err != nil {
			s.logger.Errorf("Error closing listener: %v", err)
		}

		// Wait for active connections to finish processing
		done := make(chan struct{})
		go func() {
			s.wg.Wait()
			close(done)
		}()

		// Use shutdown timeout to avoid hanging indefinitely
		select {
		case <-done:
			s.logger.Info("All connections closed. Server stopped.")
		case <-time.After(s.config.ShutdownTimeout):
			s.logger.Error("Timeout reached. Forcing shutdown.")
		}

		// Cancel context to free up resources and signal shutdown
		s.cancel()
	})
}
