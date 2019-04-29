package handlers

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/denouche/plex-watcher/middlewares"
	"github.com/denouche/plex-watcher/storage/dao"
	"github.com/denouche/plex-watcher/storage/dao/inmemory"
	"github.com/denouche/plex-watcher/storage/services/notifications"
	"github.com/denouche/plex-watcher/storage/services/plex"
	"github.com/denouche/plex-watcher/storage/validators"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
)

type Config struct {
	LogLevel       string
	LogFormat      string
	Port           int
	SlackURI       string
	DBInMemoryFile string
}

type handlersContext struct {
	db        dao.Database
	plex      plex.Plex
	validator *validator.Validate
	notif     notifications.Notifications
}

func NewRouter(config *Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.HandleMethodNotAllowed = true

	router.Use(gin.Recovery())
	router.Use(middlewares.GetLoggerMiddleware())
	router.Use(middlewares.GetHTTPLoggerMiddleware())

	hc := &handlersContext{}
	hc.db = inmemory.NewDatabaseInMemory(config.DBInMemoryFile)
	hc.notif = notifications.NewSlack(config.SlackURI)
	hc.plex = plex.NewPlex()
	hc.validator = newValidator()

	public := router.Group("/")
	public.Handle(http.MethodGet, "/_health", hc.GetHealth)

	public.Handle(http.MethodPost, "/scan", hc.ScanLibraries)
	public.Handle(http.MethodPost, "/command", hc.HandleCommand)

	return router
}

func newValidator() *validator.Validate {
	va := validator.New()

	va.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)
		if len(name) < 1 {
			return ""
		}
		return name[0]
	})

	for k, v := range validators.CustomValidators {
		if v.Validator != nil {
			va.RegisterValidationCtx(k, v.Validator)
		}
	}

	return va
}
