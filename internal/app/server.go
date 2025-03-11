package app

import (
	"context"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"net"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"sync"
	"syscall"
	"time"
	"word-of-wisdom/internal/config"
)

const (
	MsgOnManyReq     = "Too many requests. Please try again later.\n"
	MsgOnErrInternal = "Internal server error. Please try again later.\n"
)

// Server encapsulates the TCP server's behavior
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
	limiterMap   sync.Map
}

// NewServer initializes a new server instance
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
		s.logger.Fatalf("Failed to start server: %v", err)
		return
	}

	s.logger.Infof("Server started on port %s", s.config.Port)

	go s.acceptConnections()

	// Wait for shutdown signal
	<-s.ctx.Done()
	s.Shutdown()
}

// acceptConnections listens for incoming connections and limits concurrency
func (s *Server) acceptConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.ctx.Err() != nil {
				s.logger.Info("Server is shutting down, stopping connection handling...")
				return
			}
			if strings.Contains(err.Error(), "use of closed network connection") {
				s.logger.Info("Listener closed, stopping connection handling...")
				return
			}
			s.logger.Errorf("Failed to accept connection: %v", err)
			continue
		}

		select {
		case s.semaphore <- struct{}{}:
			s.wg.Add(1)
			go s.handleClient(conn)
		default:
			s.logger.Warn("Too many connections. Rejecting client.")
			_ = conn.Close()
		}
	}
}

// getLimiterForIP returns a rate limiter per IP
func (s *Server) getLimiterForIP(ip string) *rate.Limiter {
	limiter, loaded := s.limiterMap.LoadOrStore(ip, rate.NewLimiter(rate.Every(100*time.Millisecond), s.config.RateLimitEvery100MS))
	if !loaded {
		s.logger.Infof("Created new rate limiter for IP: %s", ip)
	}
	return limiter.(*rate.Limiter)
}

// handleClient processes a single client connection
func (s *Server) handleClient(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()
	defer func() { <-s.semaphore }() // Release slot
	defer s.recoverPanic("handleClient", conn)

	ip := conn.RemoteAddr().(*net.TCPAddr).IP.String()
	limiter := s.getLimiterForIP(ip)

	if !limiter.Allow() {
		_, _ = conn.Write([]byte(MsgOnManyReq))
		return
	}

	if err := conn.SetDeadline(time.Now().Add(s.config.ConnectionTimeout)); err != nil {
		s.logger.Errorf("Failed to set deadline for client %s: %v", ip, err)
	}

	if err := s.handler.HandleConnection(conn); err != nil {
		s.logger.Errorf("Error handling client %s: %v", ip, err)
	}
}

// recoverPanic handles panics and logs stack traces
func (s *Server) recoverPanic(funcName string, conn net.Conn) {
	if r := recover(); r != nil {
		s.logger.Errorf("Panic recovered in %s: %v\nStack trace:\n%s", funcName, r, string(debug.Stack()))
		if conn != nil {
			_, _ = conn.Write([]byte(MsgOnErrInternal))
		}
	}
}

// Shutdown gracefully stops the server
func (s *Server) Shutdown() {
	s.shutdownOnce.Do(func() {
		s.logger.Info("Shutting down server...")

		if err := s.listener.Close(); err != nil {
			s.logger.Errorf("Error closing listener: %v", err)
		}

		done := make(chan struct{})
		go func() {
			s.wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			s.logger.Info("All connections closed. Server stopped.")
		case <-time.After(s.config.ShutdownTimeout):
			s.logger.Warn("Shutdown timeout reached. Forcing termination.")
		}

		s.cancel()
	})
}
