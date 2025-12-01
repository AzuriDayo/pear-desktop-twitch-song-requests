package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/appservices"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/helpers"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	helpers.PreflightTest()

	app := NewApp()
	log.Fatalln(app.Run())
}

type App struct {
	ctx           context.Context
	twitchService appservices.TwitchWS
}

func NewApp() *App {
	return &App{
		ctx: context.TODO(),
	}
}

//go:embed build/*
var staticControlPanelFS embed.FS

func (a *App) Run() error {
	log.Println("App is running on port 3999...")
	http.HandleFunc("GET /api/twitch-oauth", a.handleTwitchOAuth)
	buildFS, err := fs.Sub(staticControlPanelFS, "build")
	if err != nil {
		panic(err)
	}
	http.Handle("GET /", http.FileServer(http.FS(buildFS)))

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
	args = append(args, "http://127.0.0.1:3999/")
	exec.Command(cmd, args...).Start()
	return http.ListenAndServe(":3999", nil)
}

func (a *App) handleTwitchOAuth(w http.ResponseWriter, r *http.Request) {
	// auth data in url hash string params as get request
}
