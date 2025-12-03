package main

import (
	"context"
	"embed"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/appservices"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/helpers"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/staticservices"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	helpers.PreflightTest()

	app := NewApp()
	log.Fatalln(app.Run())
}

type App struct {
	pearDesktopService *staticservices.PearDesktopService
	pearDesktopWS      *appservices.PearDesktopWS
	ctx                context.Context
	cancel             context.CancelFunc
}

func NewApp() *App {
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize Pear Desktop static service
	pearDesktopService := staticservices.NewPearDesktopService()

	// Initialize Pear Desktop websocket service
	pearDesktopWS := appservices.NewPearDesktopWS()

	return &App{
		pearDesktopService: pearDesktopService,
		pearDesktopWS:      pearDesktopWS,
		ctx:                ctx,
		cancel:             cancel,
	}
}

//go:embed build/*
var staticControlPanelFS embed.FS

func (a *App) Run() error {
	log.Println("App is running on port 3999...")

	// Ensure services are stopped on exit
	defer a.cancel()

	// Test Pear Desktop connection
	if err := a.pearDesktopService.TestConnection(); err != nil {
		log.Printf("Warning: Pear Desktop service not available: %v", err)
	} else {
		log.Println("Pear Desktop service connected successfully")
	}

	// Start Pear Desktop websocket service
	if err := a.pearDesktopWS.StartCtx(a.ctx); err != nil {
		log.Printf("Failed to start Pear Desktop WS service: %v", err)
		return err
	}

	// Start goroutine to handle music player state updates
	go a.handleMusicPlayerUpdates()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())
	e.StaticFS("/", echo.MustSubFS(staticControlPanelFS, "build"))

	apiV1 := e.Group("/api/v1")
	apiV1.POST("/twitch-oauth", a.handleTwitchOAuth)
	apiV1.GET("/music/state", a.handleGetMusicState)
	apiV1.POST("/music/state", a.handleSetMusicState)

	var cmd string
	var args []string
	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, "http://localhost:3999/") // must use localhost here because twitch does not allow 127.0.0.1
	exec.Command(cmd, args...).Start()

	// Start server with graceful shutdown
	go func() {
		<-a.ctx.Done()
		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	return e.Start("127.0.0.1:3999")
}

func (a *App) handleTwitchOAuth(c echo.Context) error {
	// auth data in url hash string params as get request
	body := c.Request().Body
	rawBodyData, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	defer body.Close()

	authData := struct {
		Something string
	}{}
	err = json.Unmarshal(rawBodyData, &authData)
	return echo.NewHTTPError(http.StatusTeapot, "under construction", err)
}

func (a *App) handleGetMusicState(c echo.Context) error {
	state, err := a.pearDesktopService.GetMusicPlayerState()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get music player state", err)
	}
	return c.JSON(http.StatusOK, state)
}

func (a *App) handleSetMusicState(c echo.Context) error {
	var state staticservices.MusicPlayerState
	if err := c.Bind(&state); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body", err)
	}

	if err := a.pearDesktopService.SetMusicPlayerState(&state); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to set music player state", err)
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (a *App) handleMusicPlayerUpdates() {
	for {
		select {
		case <-a.ctx.Done():
			return
		case update := <-a.pearDesktopWS.RcvChan():
			log.Printf("Received music player state update: %+v", update)
			// Here you can add logic to handle state updates, such as:
			// - Update database
			// - Broadcast to other services
			// - Trigger events based on state changes
		}
	}
}
