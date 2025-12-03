package appservices

import (
	"context"
	"encoding/json"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/coder/websocket"
)

// MusicPlayerStateUpdate represents a music player state update from the websocket
type MusicPlayerStateUpdate struct {
	IsPlaying   bool      `json:"isPlaying"`
	CurrentSong string    `json:"currentSong,omitempty"`
	Volume      int       `json:"volume,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
}

// PearDesktopWS handles websocket connection to Pear Desktop for music player state updates
type PearDesktopWS struct {
	stopChan          chan struct{}
	wsURL             string
	conn              *websocket.Conn
	log               *log.Logger
	msgChan           chan MusicPlayerStateUpdate
	rcvChan           chan MusicPlayerStateUpdate
	reconnectInterval time.Duration
}

// StartCtx starts the websocket connection and begins listening for state updates
func (s *PearDesktopWS) StartCtx(ctx context.Context) error {
	s.log.Println("Pear Desktop WS service starting...")

	// Initial connection attempt
	if err := s.connect(); err != nil {
		return err
	}

	// Start message handling goroutine
	go s.handleMessages()

	// Start reconnection handler
	go s.handleReconnection(ctx)

	// Stop handler
	go func() {
		<-ctx.Done()
		s.Stop()
	}()

	s.log.Println("Pear Desktop WS service started.")
	return nil
}

// connect establishes the websocket connection
func (s *PearDesktopWS) connect() error {
	u, err := url.Parse(s.wsURL)
	if err != nil {
		return err
	}

	conn, _, err := websocket.Dial(context.Background(), u.String(), nil)
	if err != nil {
		return err
	}

	s.conn = conn
	return nil
}

// handleMessages processes incoming websocket messages
func (s *PearDesktopWS) handleMessages() {
	defer func() {
		if r := recover(); r != nil {
			s.log.Println("Recovered in handleMessages:", r)
		}
	}()

	for {
		select {
		case <-s.stopChan:
			return
		default:
			if s.conn == nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			_, message, err := s.conn.Read(context.Background())
			if err != nil {
				s.log.Printf("WebSocket read error: %v", err)
				s.conn = nil
				continue
			}

			var update MusicPlayerStateUpdate
			if err := json.Unmarshal(message, &update); err != nil {
				s.log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			// Send to receive channel for external listeners
			select {
			case s.rcvChan <- update:
			default:
				s.log.Println("Receive channel full, dropping message")
			}
		}
	}
}

// handleReconnection manages reconnection logic
func (s *PearDesktopWS) handleReconnection(ctx context.Context) {
	ticker := time.NewTicker(s.reconnectInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stopChan:
			return
		case <-ticker.C:
			if s.conn == nil {
				s.log.Println("Attempting to reconnect to Pear Desktop WS...")
				if err := s.connect(); err != nil {
					s.log.Printf("Reconnection failed: %v", err)
				} else {
					s.log.Println("Reconnected to Pear Desktop WS")
				}
			}
		}
	}
}

// Stop closes the websocket connection and stops the service
func (s *PearDesktopWS) Stop() error {
	defer func() {
		if r := recover(); r != nil {
			s.log.Println("Recovered in PearDesktopWS Stop():", r)
		}
	}()

	s.log.Println("Pear Desktop WS service stopping...")
	close(s.stopChan)

	if s.conn != nil {
		err := s.conn.Close(websocket.StatusNormalClosure, "service stopping")
		if err != nil {
			s.log.Printf("Error closing websocket connection: %v", err)
		}
		s.conn = nil
	}

	s.log.Println("Pear Desktop WS service stopped.")
	return nil
}

// MsgChan returns the message channel (for sending messages if needed)
func (s *PearDesktopWS) MsgChan() chan MusicPlayerStateUpdate {
	return s.msgChan
}

// RcvChan returns the receive channel for incoming music player state updates
func (s *PearDesktopWS) RcvChan() chan MusicPlayerStateUpdate {
	return s.rcvChan
}

// Log returns the logger
func (s *PearDesktopWS) Log() *log.Logger {
	return s.log
}

// NewPearDesktopWS creates a new Pear Desktop websocket service
func NewPearDesktopWS() *PearDesktopWS {
	stopChan := make(chan struct{})
	return &PearDesktopWS{
		stopChan:          stopChan,
		wsURL:             "ws://localhost:26538/ws",
		log:               log.New(os.Stderr, "PEAR_DESKTOP_WS ", log.Ldate|log.Ltime),
		msgChan:           make(chan MusicPlayerStateUpdate, 100),
		rcvChan:           make(chan MusicPlayerStateUpdate, 100),
		reconnectInterval: 5 * time.Second,
	}
}
