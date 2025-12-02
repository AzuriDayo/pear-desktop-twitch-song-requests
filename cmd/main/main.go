package main

import (
	"embed"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/helpers"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var DEBUG string

func main() {
	if DEBUG == "true" {
		if m, err := godotenv.Read(".env"); err == nil {
			err = os.Setenv("VITE_TWITCH_CLIENT_ID", m["VITE_TWITCH_CLIENT_ID"])
			if err != nil {
				panic(err)
			}
			log.Println("DEBUG: Overriding VITE_TWITCH_CLIENT_ID:", m["VITE_TWITCH_CLIENT_ID"])
		}
	}
	helpers.PreflightTest()

	app := NewApp()
	log.Fatalln(app.Run())
}

type App struct {
	// twitchService appservices.TwitchWS
}

func NewApp() *App {
	return &App{}
}

//go:embed build/*
var staticControlPanelFS embed.FS

func (a *App) Run() error {
	log.Println("App is running on port 3999...")
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())
	e.StaticFS("/", echo.MustSubFS(staticControlPanelFS, "build"))

	apiV1 := e.Group("/api/v1")
	apiV1.POST("/twitch-oauth", a.handleTwitchOAuth)

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
