package main

import (
	"context"
	"embed"
	"log"
	"time"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/appservices"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/data"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/helpers"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/songrequests"
	"github.com/nicklaw5/helix/v2"
	"github.com/recws-org/recws"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

var version = "development"

type twitchData struct {
	accessToken            string
	refreshToken           string
	refreshTokenLastUsedAt time.Time
	login                  string
	userID                 string
	isAuthenticated        bool
	expiresDate            time.Time
}

//go:embed all:frontend
var assets embed.FS

func main() {
	acquireSingleInstanceLock()
	defer releaseSingleInstanceLock()
	setTitle("Pear Desktop Twitch Song Requests by AzuriDayo_")
	log.Println("Starting Pear Desktop Twitch Song Requests", version)
	go checkForUpdates()

	helpers.PreflightTest()
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "Pear Desktop Twitch Song Requests",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:  app.startup,
		OnShutdown: app.shutdown,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
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
	wailsCtx                context.Context
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
		pearDesktopIncomingMsgs: make(chan []byte),
		cmdPermissions:          defaultCmdPermissions(),
	}
}

// startup is invoked by Wails once the window is created. It stores the Wails
// context (needed for runtime.EventsEmit) and kicks off all background services.
func (a *App) startup(ctx context.Context) {
	a.wailsCtx = ctx
	go func() {
		if err := a.runBackground(); err != nil {
			log.Println("background services error:", err)
		}
	}()
}

// shutdown is invoked by Wails when the application is about to quit.
func (a *App) shutdown(ctx context.Context) {
	a.cancel()
}

// emit is a thin wrapper around wails runtime EventsEmit that is safe to call
// before the Wails context is ready (events are simply dropped in that case).
func (a *App) emit(eventName string, optionalData ...interface{}) {
	if a.wailsCtx == nil {
		return
	}
	wailsruntime.EventsEmit(a.wailsCtx, eventName, optionalData...)
}

// runBackground starts every long-running service: SQLite load, the Pear Desktop
// websocket bridge, and the Twitch EventSub clients.
func (a *App) runBackground() error {
	// load sqlite
	err := a.loadSqliteSettings()
	if err != nil {
		return err
	}

	notifyIfTokenExpired("main", a.twitchDataStruct.accessToken, a.twitchDataStruct.isAuthenticated)
	a.notifyImplicitGrantAccessTokenExpiresSoon(false)
	notifyRefreshTokenExpiresSoon("main", a.twitchDataStruct.refreshTokenLastUsedAt, a.twitchDataStruct.refreshToken)
	if a.twitchDataStructBot.accessToken != "" || a.twitchDataStructBot.refreshToken != "" {
		notifyIfTokenExpired("bot", a.twitchDataStructBot.accessToken, a.twitchDataStructBot.isAuthenticated)
		a.notifyImplicitGrantAccessTokenExpiresSoon(true)
		notifyRefreshTokenExpiresSoon("bot", a.twitchDataStructBot.refreshTokenLastUsedAt, a.twitchDataStructBot.refreshToken)
	}

	go a.runTwitchTokenMaintenance()

	// Auto reconnect pear desktop and funnel messages to channel
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

	twitchTokenExpiresSoon := a.twitchDataStruct.isAuthenticated && time.Now().Add(-15*24*time.Hour).After(a.twitchDataStruct.expiresDate)
	if a.twitchDataStruct.isAuthenticated && twitchTokenExpiresSoon {
		log.Println("ALERT! Main account Token expiry is soon, consider refreshing token.")
	}
	twitchTokenBotExpiresSoon := a.twitchDataStructBot.isAuthenticated && time.Now().Add(-15*24*time.Hour).After(a.twitchDataStructBot.expiresDate)
	if a.twitchDataStructBot.isAuthenticated && twitchTokenBotExpiresSoon {
		log.Println("ALERT! Bot Token expiry is soon, consider refreshing token.")
	}

	return nil
}
