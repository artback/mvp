package graceful

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	*http.Server
	shutdownFinished chan struct{}
}

func (s *Server) ListenAndServe() (err error) {
	if s.shutdownFinished == nil {
		s.shutdownFinished = make(chan struct{})
	}

	err = s.Server.ListenAndServe()
	if err == http.ErrServerClosed {
		// expected error after calling Server.Shutdown().
		err = nil
	} else if err != nil {
		return fmt.Errorf("unexpected error from ListenAndServe: %w", err)
	}

	log.Println("waiting for shutdown finishing...")
	<-s.shutdownFinished
	log.Println("shutdown finished")

	return
}

func (s *Server) WaitForExitingSignal(timeout time.Duration) {
	var waiter = make(chan os.Signal, 1) // buffered channel
	signal.Notify(waiter, syscall.SIGTERM, syscall.SIGINT)

	// blocks here until there's a signal
	<-waiter

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := s.Server.Shutdown(ctx)
	if err != nil {
		log.Println("shutting down: " + err.Error())
	} else {
		log.Println("shutdown processed successfully")
		close(s.shutdownFinished)
	}
}
