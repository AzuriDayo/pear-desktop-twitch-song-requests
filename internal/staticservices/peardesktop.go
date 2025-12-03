package staticservices

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// MusicPlayerState represents the current state of the music player
type MusicPlayerState struct {
	IsPlaying  bool   `json:"isPlaying"`
	CurrentSong string `json:"currentSong,omitempty"`
	Volume     int    `json:"volume,omitempty"`
}

// PearDesktopService handles communication with the Pear Desktop background process
type PearDesktopService struct {
	baseURL    string
	httpClient *http.Client
}

// NewPearDesktopService creates a new Pear Desktop service instance
func NewPearDesktopService() *PearDesktopService {
	return &PearDesktopService{
		baseURL: "http://localhost:26538",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// TestConnection verifies the connection to the Pear Desktop background process
func (s *PearDesktopService) TestConnection() error {
	resp, err := s.httpClient.Get(s.baseURL + "/health")
	if err != nil {
		return fmt.Errorf("failed to connect to Pear Desktop: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Pear Desktop health check failed with status: %d", resp.StatusCode)
	}

	return nil
}

// GetMusicPlayerState retrieves the current music player state
func (s *PearDesktopService) GetMusicPlayerState() (*MusicPlayerState, error) {
	resp, err := s.httpClient.Get(s.baseURL + "/api/music/state")
	if err != nil {
		return nil, fmt.Errorf("failed to get music player state: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get music player state, status: %d", resp.StatusCode)
	}

	var state MusicPlayerState
	if err := json.NewDecoder(resp.Body).Decode(&state); err != nil {
		return nil, fmt.Errorf("failed to decode music player state: %w", err)
	}

	return &state, nil
}

// SetMusicPlayerState updates the music player state
func (s *PearDesktopService) SetMusicPlayerState(state *MusicPlayerState) error {
	jsonData, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	resp, err := s.httpClient.Post(s.baseURL+"/api/music/state", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to set music player state: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to set music player state, status: %d", resp.StatusCode)
	}

	return nil
}