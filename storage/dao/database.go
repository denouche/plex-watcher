package dao

import (
	"github.com/denouche/plex-watcher/storage/model"
	"github.com/gin-gonic/gin"
)

type Database interface {
	GetAllLibraries() ([]*model.Library, error)
	GetLibrariesByTeamIDAndUserID(teamID, userID string) ([]*model.Library, error)
	AddLibrary(l *model.Library) error
	SaveLibrary(c *gin.Context, l *model.Library) error
}
