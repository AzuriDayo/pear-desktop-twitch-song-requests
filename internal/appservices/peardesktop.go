package appservices

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/data"
	"github.com/coder/websocket"
)

// SongInfo represents detailed song information from Pear Desktop websocket
type SongInfo struct {
	Title            string    `json:"title"`
	AlternativeTitle string    `json:"alternativeTitle,omitempty"`
	Artist           string    `json:"artist"`
	ArtistURL        string    `json:"artistUrl,omitempty"`
	Views            int       `json:"views,omitempty"`
	UploadDate       string    `json:"uploadDate,omitempty"`
	URL              string    `json:"url"`
	SongDuration     int       `json:"songDuration"`
	ImageSrc         string    `json:"imageSrc"`
	Image            ImageInfo `json:"image,omitempty"`
	ElapsedSeconds   int       `json:"elapsedSeconds"`
	IsPaused         bool      `json:"isPaused"`
	Album            string    `json:"album,omitempty"`
	VideoID          string    `json:"videoId,omitempty"`
	PlaylistID       string    `json:"playlistId,omitempty"`
	MediaType        string    `json:"mediaType,omitempty"`
	Tags             []string  `json:"tags,omitempty"`
}

// ImageInfo represents image information
type ImageInfo struct {
	IsMacTemplateImage bool `json:"isMacTemplateImage"`
}

// PositionChangedMessage represents a position update from the websocket
type PositionChangedMessage struct {
	Type     string `json:"type"` // "POSITION_CHANGED"
	Position int    `json:"position"`
}

// VideoChangedMessage represents a video change from the websocket
type VideoChangedMessage struct {
	Type     string   `json:"type"` // "VIDEO_CHANGED"
	Song     SongInfo `json:"song"`
	Position int      `json:"position"`
}

// PlayerInfoMessage represents a PLAYER_INFO message from Pear Desktop
type PlayerInfoMessage struct {
	Type      string   `json:"type"` // "PLAYER_INFO"
	Song      SongInfo `json:"song"`
	IsPlaying bool     `json:"isPlaying"`
	Muted     bool     `json:"muted"`
	Position  int      `json:"position"`
	Volume    int      `json:"volume"`
	Repeat    string   `json:"repeat"`
	Shuffle   bool     `json:"shuffle"`
}

// WebSocketMessage represents any websocket message from Pear Desktop
type WebSocketMessage struct {
	Type      string   `json:"type"`
	Position  int      `json:"position,omitempty"`
	IsPlaying bool     `json:"isPlaying"`
	Song      SongInfo `json:"song"`
}

// WebSocketStateUpdate represents a music player state update from the websocket
type WebSocketStateUpdate struct {
	IsPlaying      *bool     `json:"isPlaying"`
	CurrentSong    string    `json:"currentSong,omitempty"`
	Artist         string    `json:"artist,omitempty"`
	URL            string    `json:"url,omitempty"`
	SongDuration   int       `json:"songDuration,omitempty"`
	ImageSrc       string    `json:"imageSrc,omitempty"`
	ElapsedSeconds int       `json:"elapsedSeconds,omitempty"`
	Timestamp      time.Time `json:"timestamp"`
}

// PearDesktopService handles websocket connection to Pear Desktop for music player state updates
type PearDesktopService struct {
	stopChan          chan struct{}
	wsURL             string
	conn              *websocket.Conn
	log               *log.Logger
	msgChan           chan WebSocketStateUpdate
	rcvChan           chan WebSocketStateUpdate
	reconnectInterval time.Duration
}

// StartCtx starts the websocket connection and begins listening for state updates
func (s *PearDesktopService) StartCtx(ctx context.Context) error {
	s.log.Println("Pear Desktop WS service starting...")

	// Start reconnection handler (this will handle initial connection and retries)
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
func (s *PearDesktopService) connect() error {
	u, err := url.Parse(s.wsURL)
	if err != nil {
		return err
	}

	// Create headers with Authorization (Bearer token can be empty)
	headers := http.Header{}
	headers.Set("Authorization", "Bearer ")

	s.log.Printf("Attempting websocket connection to: %s", u.String())
	conn, _, err := websocket.Dial(context.Background(), u.String(), &websocket.DialOptions{
		HTTPHeader: headers,
	})
	if err != nil {
		s.log.Printf("WebSocket connection failed: %v", err)
		return err
	}

	s.log.Println("WebSocket connection established successfully")
	s.conn = conn
	return nil
}

// handleMessages processes incoming websocket messages
func (s *PearDesktopService) handleMessages() {
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

			// First, determine the message type by parsing just the type field
			var typeCheck struct {
				Type string `json:"type"`
			}
			if err := json.Unmarshal(message, &typeCheck); err != nil {
				s.log.Printf("Failed to unmarshal message type: %v", err)
				continue
			}

			// Create state update based on message type
			var update WebSocketStateUpdate
			update.Timestamp = time.Now()

			switch typeCheck.Type {
			case "PLAYER_INFO":
				var playerMsg PlayerInfoMessage
				if err := json.Unmarshal(message, &playerMsg); err != nil {
					s.log.Printf("Failed to unmarshal PLAYER_INFO message: %v", err)
					continue
				}
				// s.log.Println("Received PLAYER_INFO:", string(message))

				// Handle PLAYER_INFO message
				if playerMsg.Song.Title != "" {
					update.CurrentSong = playerMsg.Song.Title
					update.Artist = playerMsg.Song.Artist
					update.URL = playerMsg.Song.URL
					update.SongDuration = playerMsg.Song.SongDuration
					update.ImageSrc = playerMsg.Song.ImageSrc
					update.ElapsedSeconds = playerMsg.Song.ElapsedSeconds
					update.IsPlaying = &playerMsg.IsPlaying

					s.log.Printf("Player info - Title: %s, Artist: %s, Duration: %ds, Position: %ds, Playing: %t",
						playerMsg.Song.Title, playerMsg.Song.Artist, playerMsg.Song.SongDuration, playerMsg.Song.ElapsedSeconds, playerMsg.IsPlaying)
				}

			case "POSITION_CHANGED":
				var posMsg PositionChangedMessage
				if err := json.Unmarshal(message, &posMsg); err != nil {
					s.log.Printf("Failed to unmarshal POSITION_CHANGED message: %v", err)
					continue
				}
				// Position changed - update elapsed seconds
				update.ElapsedSeconds = posMsg.Position

			case "PLAYER_STATE_CHANGED", "VIDEO_CHANGED":
				var wsMsg WebSocketMessage
				if err := json.Unmarshal(message, &wsMsg); err != nil {
					s.log.Printf("Failed to unmarshal websocket message: %v", err)
					continue
				}
				// s.log.Println("Received", wsMsg.Type, ":", string(message))

				switch wsMsg.Type {
				case "PLAYER_STATE_CHANGED":
					update.IsPlaying = &wsMsg.IsPlaying
					update.ElapsedSeconds = wsMsg.Position
				case "VIDEO_CHANGED":
					// Video changed - extract song information
					if wsMsg.Song.Title != "" {
						update.CurrentSong = wsMsg.Song.Title
						update.Artist = wsMsg.Song.Artist
						update.URL = wsMsg.Song.URL
						update.SongDuration = wsMsg.Song.SongDuration
						update.ImageSrc = wsMsg.Song.ImageSrc
						update.ElapsedSeconds = wsMsg.Position
						b := !wsMsg.Song.IsPaused
						update.IsPlaying = &b

						s.log.Printf("Video changed - Title: %s, Artist: %s, Duration: %d",
							wsMsg.Song.Title, wsMsg.Song.Artist, wsMsg.Song.SongDuration)
					}
				}

			default:
				s.log.Printf("Unknown message type: %s", typeCheck.Type)
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
func (s *PearDesktopService) handleReconnection(ctx context.Context) {
	// Initial connection attempt
	s.log.Println("Attempting initial connection to Pear Desktop WS...")
	if err := s.connect(); err != nil {
		s.log.Printf("Initial connection failed: %v", err)
	} else {
		s.log.Println("Connected to Pear Desktop WS")
		// Start message handling goroutine when connected
		go s.handleMessages()
	}

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
					// Start message handling goroutine when reconnected
					go s.handleMessages()
				}
			}
		}
	}
}

// Stop closes the websocket connection and stops the service
func (s *PearDesktopService) Stop() error {
	defer func() {
		if r := recover(); r != nil {
			s.log.Println("Recovered in PearDesktopService Stop():", r)
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
func (s *PearDesktopService) MsgChan() chan WebSocketStateUpdate {
	return s.msgChan
}

// RcvChan returns the receive channel for incoming music player state updates
func (s *PearDesktopService) RcvChan() chan WebSocketStateUpdate {
	return s.rcvChan
}

// Log returns the logger
func (s *PearDesktopService) Log() *log.Logger {
	return s.log
}

// NewPearDesktopService creates a new Pear Desktop websocket service
func NewPearDesktopService() *PearDesktopService {
	stopChan := make(chan struct{})
	return &PearDesktopService{
		stopChan:          stopChan,
		wsURL:             "ws://" + data.GetPearDesktopHost() + "/api/v1/ws",
		log:               log.New(os.Stderr, "PEAR_DESKTOP_WS ", log.Ldate|log.Ltime),
		msgChan:           make(chan WebSocketStateUpdate, 100),
		rcvChan:           make(chan WebSocketStateUpdate, 100),
		reconnectInterval: 5 * time.Second,
	}
}
