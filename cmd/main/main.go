package main

import (
	"bufio"
	"context"
	"embed"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/appservices"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/data"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/helpers"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/songrequests"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nicklaw5/helix/v2"
	"github.com/recws-org/recws"
	"golang.org/x/net/websocket"
)

var version = "development"

type twitchData struct {
	accessToken     string
	login           string
	userID          string
	isAuthenticated bool
	expiresDate     time.Time
}

func main() {
	acquireSingleInstanceLock()
	defer releaseSingleInstanceLock()
	setTitle("Pear Desktop Twitch Song Requests by AzuriDayo_")
	log.Println("Starting Pear Desktop Twitch Song Requests", version)
	go checkForUpdates()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	helpers.PreflightTest()
	app := NewApp()

	go func() {
		log.Println(app.Run())
	}()
	<-sigs
	app.cancel()

	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

type App struct {
	twitchDataStruct        *twitchData
	twitchDataStructBot     *twitchData
	helix                   *helix.Client
	helixBot                *helix.Client
	twitchWSService         *appservices.TwitchWS
	twitchWSBotService      *appservices.TwitchWS
	streamOnline            bool
	pearDesktopIncomingMsgs chan []byte
	ctx                     context.Context
	cancel                  context.CancelFunc
	clients                 map[*websocket.Conn]struct{}
	clientsMu               sync.RWMutex
	clientsBroadcast        chan string
	songRequestRewardID     string
	cmdPermissions          map[string]int
}

// defaultCmdPermissions returns the default minimum permission levels for each command.
func defaultCmdPermissions() map[string]int {
	return map[string]int{
		data.DB_KEY_CMD_PERMISSION_SR:      data.PermissionLevelSubscriber,
		data.DB_KEY_CMD_PERMISSION_QUEUE:   data.PermissionLevelViewer,
		data.DB_KEY_CMD_PERMISSION_SONG:    data.PermissionLevelViewer,
		data.DB_KEY_CMD_PERMISSION_DELSONG: data.PermissionLevelModerator,
	}
}

func NewApp() *App {
	ctx, cancel := context.WithCancel(context.Background())
	c, err := helix.NewClient(&helix.Options{
		ClientID: data.GetTwitchClientID(),
	})
	if err != nil {
		log.Fatal("Failed to create helix client: ", err)
	}
	c2, err := helix.NewClient(&helix.Options{
		ClientID: data.GetTwitchClientID(),
	})
	if err != nil {
		log.Fatal("Failed to create helix bot client: ", err)
	}
	return &App{
		twitchDataStruct:        &twitchData{},
		twitchDataStructBot:     &twitchData{},
		ctx:                     ctx,
		cancel:                  cancel,
		helix:                   c,
		helixBot:                c2,
		clientsBroadcast:        make(chan string),
		clientsMu:               sync.RWMutex{},
		clients:                 make(map[*websocket.Conn]struct{}),
		pearDesktopIncomingMsgs: make(chan []byte),
		cmdPermissions:          defaultCmdPermissions(),
	}
}

//go:embed build/*
var staticControlPanelFS embed.FS

func (a *App) Run() error {
	// load sqlite
	err := a.loadSqliteSettings()
	if err != nil {
		return err
	}

	notifyIfTokenExpired("main", a.twitchDataStruct.accessToken, a.twitchDataStruct.isAuthenticated)
	notifyTokenExpiresSoon("main", a.twitchDataStruct.expiresDate, a.twitchDataStruct.isAuthenticated)
	if a.twitchDataStructBot.accessToken != "" {
		notifyIfTokenExpired("bot", a.twitchDataStructBot.accessToken, a.twitchDataStructBot.isAuthenticated)
		notifyTokenExpiresSoon("bot", a.twitchDataStructBot.expiresDate, a.twitchDataStructBot.isAuthenticated)
	}

	// Auto reconnect pear desktop and funnel mesasges to channel
	log.Println("Pear Desktop WS service starting...")
	ws := recws.RecConn{
		RecIntvlFactor: 1,               // multiplier backoff
		RecIntvlMin:    3 * time.Second, // start time
		NonVerbose:     true,
		SubscribeHandler: func() error {
			log.Println("Connected to Pear Desktop")
			return nil
		},
	}
	ws.Dial("ws://"+songrequests.GetPearDesktopHost()+"/api/v1/ws", nil)
	go func() {
		for {
			select {
			case <-a.ctx.Done():
				go ws.Close()
				return
			default:
				if !ws.IsConnected() {
					time.Sleep(3 * time.Second)
					continue
				}
				_, message, err := ws.Conn.ReadMessage()
				if err != nil {
					time.Sleep(3 * time.Second)
					continue
				}

				a.pearDesktopIncomingMsgs <- message
			}

		}
	}()

	// Handle Pear desktop messages
	go a.handlePearDesktopMsgs()

	// Auto reconnect twitch ws
	go func() {
		twitchWSService := appservices.NewTwitchWS(a.helix, &a.twitchDataStruct.userID, &a.twitchDataStruct.login, nil, nil, nil, songrequests.GetSubscriptions(), a.SetSubscriptionHandlers, false)
		a.twitchWSService = twitchWSService
		for {
			if a.helix.GetUserAccessToken() != "" {
				valid, _, _ := a.helix.ValidateToken(a.helix.GetUserAccessToken())
				if valid {
					err := a.twitchWSService.StartCtx(a.ctx)
					if err == nil {
						// graceful shutdown
						return
					}
					log.Println("Twitch WS MAIN disconnected, attempt to reconnect")
				}
				// always sleep 5s after token validation
				time.Sleep(5 * time.Second)
			} else {
				time.Sleep(5 * time.Second)
			}
		}
	}()

	// Auto reconnect twitch ws bot
	go func() {
		twitchWSBotService := appservices.NewTwitchWS(a.helixBot, &a.twitchDataStructBot.userID, &a.twitchDataStructBot.login, a.helix, &a.twitchDataStruct.userID, &a.twitchDataStruct.login, songrequests.GetSubscriptionsBot(), a.SetSubscriptionHandlersBot, true)
		a.twitchWSBotService = twitchWSBotService
		for {
			if a.helixBot.GetUserAccessToken() != "" {
				valid, _, _ := a.helixBot.ValidateToken(a.helixBot.GetUserAccessToken())
				if valid {
					err := a.twitchWSBotService.StartCtx(a.ctx)
					if err == nil {
						// graceful shutdown
						return
					}
					log.Println("Twitch WS BOT disconnected, attempt to reconnect")
				}
				// always sleep 5s after token validation
				time.Sleep(5 * time.Second)
			} else {
				time.Sleep(5 * time.Second)
			}
		}
	}()

	// Send msgs to ws clients
	go func() {
		for {
			data := <-a.clientsBroadcast
			a.clientsMu.Lock()
			for ws := range a.clients {
				websocket.Message.Send(ws, data)
			}
			a.clientsMu.Unlock()
		}
	}()

	// Echo instance
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       "build",
		Index:      "index.html",
		Filesystem: http.FS(staticControlPanelFS),
		HTML5:      true,
	}))

	apiV1 := e.Group("/api/v1")
	apiV1.POST("/twitch-oauth", a.handleApiV1TwitchOAuthPOST)
	apiV1.GET("/settings", a.handleApiV1SettingsGET)
	apiV1.PATCH("/settings", a.handleApiV1SettingsPATCH)
	apiV1.GET("/ws", a.handleApiV1WsGET)
	apiV1.DELETE("/queue/:idx", a.handleApiV1QueueDeleteDELETE)
	apiV1Requesters := apiV1.Group("/requesters")
	apiV1Requesters.GET("/history", a.handleApiV1RequestersHistoryGET)
	apiV1Twitch := apiV1.Group("/twitch")
	apiV1Twitch.GET("/custom-rewards", a.handleApiV1TwitchCustomRewardsGET)

	port := findAvailablePort()
	controlPanelURL := fmt.Sprintf("http://localhost:%d/", port)

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
	args = append(args, controlPanelURL) // must use localhost here because twitch does not allow 127.0.0.1
	twitchTokenExpiresSoon := a.twitchDataStruct.isAuthenticated && time.Now().Add(-15*24*time.Hour).After(a.twitchDataStruct.expiresDate)
	if a.twitchDataStruct.isAuthenticated && twitchTokenExpiresSoon {
		log.Println("ALERT! Main account Token expiry is soon, consider refreshing token.")
	}
	twitchTokenBotExpiresSoon := a.twitchDataStructBot.isAuthenticated && time.Now().Add(-15*24*time.Hour).After(a.twitchDataStructBot.expiresDate)
	if a.twitchDataStructBot.isAuthenticated && twitchTokenBotExpiresSoon {
		log.Println("ALERT! Bot Token expiry is soon, consider refreshing token.")
	}
	if !a.twitchDataStruct.isAuthenticated || a.songRequestRewardID == "" || twitchTokenExpiresSoon || twitchTokenBotExpiresSoon {
		exec.Command(cmd, args...).Start()
	} else {
		time.Sleep(5 * time.Second)
		log.Println("Friendly reminder, the control panel is available at " + controlPanelURL)
	}
	return e.Start(fmt.Sprintf("%s:%d", listenIP, port))
}

var listenIP = "127.0.0.1"

func isPortAvailable(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", listenIP, port))
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

// findAvailablePort tries ports 3999→3000 (decrementing), then 8080→65535 (incrementing).
func findAvailablePort() int {
	for port := 3999; port >= 3000; port-- {
		if isPortAvailable(port) {
			return port
		}
	}
	for port := 8080; port <= 65535; port++ {
		if isPortAvailable(port) {
			return port
		}
	}
	log.Fatal("No available port found")
	return 0
}
