package handlers

import (
	"net/http"

	"github.com/denouche/plex-watcher/utils"
	"github.com/gin-gonic/gin"
)

func (hc *handlersContext) GetHealth(c *gin.Context) {
	utils.JSON(c.Writer, http.StatusNoContent, nil)
}
