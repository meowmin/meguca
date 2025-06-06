package server

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/bakape/meguca/auth"
	"github.com/bakape/meguca/config"
	"github.com/bakape/meguca/db"
	"github.com/bakape/meguca/imager"
	"github.com/bakape/meguca/util"
	"github.com/bakape/meguca/websockets"
	"github.com/dimfeld/httptreemux"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/go-playground/log"

	// Add profiling to default server mux
	_ "net/http/pprof"
)

var (
	healthCheckMsg = []byte("God's in His heaven, all's right with the world")
)

// Used for overriding during tests
var webRoot = "www"

func startWebServer() (err error) {
	go func() {
		// Bind pprof to random localhost-only address
		http.ListenAndServe("localhost:0", nil)
	}()

	c := config.Server.Server

	var w strings.Builder
	w.WriteString("listening on http")
	prettyAddr := c.Address
	if len(c.Address) != 0 && c.Address[0] == ':' {
		prettyAddr = "127.0.0.1" + prettyAddr
	}
	fmt.Fprintf(&w, "://%s", prettyAddr)
	log.Info(w.String())

	gracehttp.PreStartProcess(db.Close)
	err = gracehttp.Serve(&http.Server{
		Addr:    c.Address,
		Handler: createRouter(),
	})
	if err != nil {
		return util.WrapError("error starting web server", err)
	}
	return
}

func handlePanic(w http.ResponseWriter, r *http.Request, err interface{}) {
	http.Error(w, fmt.Sprintf("500 %s", err), 500)
	ip, ipErr := auth.GetIP(r)
	if ipErr != nil {
		ip = "invalid IP"
	}
	log.Errorf("server: %s: %#v\n%s\n", ip, err, debug.Stack())
}

// Create the monolithic router for routing HTTP requests. Separated into own
// function for easier testability.
func createRouter() http.Handler {
	r := httptreemux.NewContextMux()
	r.NotFoundHandler = func(w http.ResponseWriter, _ *http.Request) {
		text404(w)
	}
	r.PanicHandler = handlePanic

	r.GET("/robots.txt", serveRobotsTXT)

	api := r.NewGroup("/api")
	api.GET("/health-check", healthCheck)
	assets := r.NewGroup("/assets")
	if config.Server.ImagerMode != config.NoImager {
		// All upload images
		api.POST("/upload", imager.NewImageUpload)
		api.POST("/upload-hash", imager.UploadImageHash)
		api.POST("/upload-megu-hash", imager.UploadMeguHash)
		api.POST("/create-thread", createThread)
		api.POST("/create-reply", createReply)

		assets.GET("/images/*path", serveImages)

		// Captcha API
		captcha := api.NewGroup("/captcha")
		captcha.GET("/:board", serveNewCaptcha)
		captcha.POST("/:board", authenticateCaptcha)
		captcha.GET("/confirmation", renderCaptchaConfirmation)
	}
	if config.Server.ImagerMode != config.ImagerOnly {
		// HTML
		r.GET("/", redirectToDefault)
		r.GET("/:board/", func(w http.ResponseWriter, r *http.Request) {
			boardHTML(w, r, extractParam(r, "board"), false)
		})
		r.GET("/:board/catalog", func(w http.ResponseWriter, r *http.Request) {
			boardHTML(w, r, extractParam(r, "board"), true)
		})
		// Needs override, because it conflicts with crossRedirect
		r.GET("/all/catalog", func(w http.ResponseWriter, r *http.Request) {
			// Artificially set board to "all"
			boardHTML(w, r, "all", true)
		})
		r.GET("/:board/:thread", threadHTML)
		r.GET("/all/:id", crossRedirect)

		html := r.NewGroup("/html")
		html.GET("/board-navigation", boardNavigation)
		html.GET("/owned-boards/:userID", ownedBoardSelection)
		html.GET("/create-board", boardCreationForm)
		html.GET("/change-password", changePasswordForm)
		html.POST("/configure-board/:board", boardConfigurationForm)
		html.POST("/configure-server", serverConfigurationForm)
		html.GET("/assign-staff/:board", staffAssignmentForm)
		html.GET("/set-banners", bannerSettingForm)
		html.GET("/set-loading", loadingAnimationForm)
		html.GET("/bans/:board", banList)
		html.GET("/mod-log/:board", modLog)
		html.GET("/report/:id", reportForm)
		html.GET("/reports/:board", reportList)

		// JSON API
		json := r.NewGroup("/json")
		boards := json.NewGroup("/boards")
		boards.GET("/:board/", func(w http.ResponseWriter, r *http.Request) {
			boardJSON(w, r, false)
		})
		boards.GET("/:board/catalog", func(w http.ResponseWriter,
			r *http.Request,
		) {
			boardJSON(w, r, true)
		})
		boards.GET("/:board/:thread", threadJSON)
		json.GET("/post/:post", servePost)
		json.GET("/config", serveConfigs)
		json.GET("/extensions", serveExtensionMap)
		json.GET("/board-config/:board", serveBoardConfigs)
		json.GET("/board-list", serveBoardList)
		json.GET("/ip-count", serveIPCount)
		json.POST("/thread-updates", serveThreadUpdates)

		// Internal API
		api.GET("/socket", func(w http.ResponseWriter, r *http.Request) {
			err := websockets.Handler(w, r)
			if err != nil {
				httpError(w, r, err)
			}
		})
		api.GET("/bitchute-title/:id", bitChuteTitle)
		api.POST("/register", register)
		api.POST("/login", login)
		api.POST("/logout", logout)
		api.POST("/logout-all", logoutAll)
		api.POST("/change-password", changePassword)
		api.POST("/board-config/:board", servePrivateBoardConfigs)
		api.POST("/configure-board/:board", configureBoard)
		api.POST("/config", servePrivateServerConfigs)
		api.POST("/configure-server", configureServer)
		api.POST("/create-board", createBoard)
		api.POST("/delete-board", deleteBoard)
		api.POST("/notification", sendNotification)
		api.POST("/assign-staff", assignStaff)
		api.POST("/same-IP/:id", getSameIPPosts)
		api.POST("/sticky", setThreadSticky)
		api.POST("/lock-thread", setThreadLock)
		api.POST("/unban/:board", unban)
		api.POST("/set-banners", setBanners)
		api.POST("/set-loading", setLoadingAnimation)
		api.POST("/report", report)
		api.GET("/sse", sse)
		api.POST("/moderate", moderate)
		api.POST("/lock-playlist", lockPlaylist)

		redir := api.NewGroup("/redirect")
		redir.POST("/by-ip", redirectByIP)
		redir.POST("/by-thread", redirectByThread)

		// Assets
		assets.GET("/banners/:board/:id", serveBanner)
		assets.GET("/loading/:board", serveLoadingAnimation)
		assets.GET("/*path", serveAssets)
	}

	return r
}

// Redirects to / requests to /all/ board
func redirectToDefault(w http.ResponseWriter, r *http.Request) {
	// Set cache-control headers to prevent caching
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")

	if config.Server.DefaultGeneralThread != nil {
		board, id, err := db.GetLatestGeneral(config.Server.DefaultGeneralThread)
		if err != nil {
			http.Redirect(w, r, "/all/", http.StatusFound)
		} else {
			http.Redirect(w, r, fmt.Sprintf("/%s/%d?last=100#bottom", board, id), http.StatusFound)
		}
	} else {
		http.Redirect(w, r, "/all/", http.StatusFound)
	}
}

// Generate a robots.txt with only select boards preventing indexing
func serveRobotsTXT(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	buf.WriteString("User-agent: *\n")

	if config.Get().GlobalDisableRobots {
		buf.WriteString("Disallow: /\n")
	} else {
		// Would be pointles without the /all/ board disallowed.
		// Also, this board can be huge. Don't want bots needlessly crawling it.
		buf.WriteString("Disallow: /all/\n")

		for _, c := range config.GetAllBoardConfigs() {
			if c.DisableRobots {
				fmt.Fprintf(&buf, "Disallow: /%s/\n", c.ID)
			}
		}
	}

	w.Header().Set("Content-Type", "text/plain")
	buf.WriteTo(w)
}

// Redirect the client to the appropriate board through a cross-board redirect
func crossRedirect(w http.ResponseWriter, r *http.Request) {
	idStr := extractParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		text404(w)
		return
	}

	board, op, err := db.GetPostParenthood(id)
	if err != nil {
		httpError(w, r, err)
		return
	}
	url := r.URL
	url.Path = fmt.Sprintf("/%s/%d", board, op)
	if url.Query().Get("last") != "" {
		url.Fragment = "bottom"
	} else {
		url.Fragment = "p" + idStr
	}
	http.Redirect(w, r, url.String(), 301)
}

// Health check to ensure server is still online
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write(healthCheckMsg)
}
