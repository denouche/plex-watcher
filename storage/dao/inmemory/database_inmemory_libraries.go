package inmemory

import (
	"encoding/json"

	"github.com/denouche/plex-watcher/storage/model"
	"github.com/denouche/plex-watcher/utils"
	"github.com/gin-gonic/gin"
	"github.com/satori/uuid"
)

const (
	cacheKeyUsers = "users"
)

func (db *DatabaseInMemory) saveLibraries(vs []*model.Library) {
	data := make([]interface{}, 0)
	for _, v := range vs {
		data = append(data, v)
	}
	db.save(cacheKeyUsers, data)
}

func (db *DatabaseInMemory) loadLibraries() []*model.Library {
	users := make([]*model.Library, 0)
	b, err := db.Cache.Get(cacheKeyUsers)
	if err != nil {
		return users
	}
	err = json.Unmarshal(b, &users)
	if err != nil {
		utils.GetLogger().Error("Error while unmarshal fake users")
	}
	return users
}

func (db *DatabaseInMemory) GetAllLibraries() ([]*model.Library, error) {
	return db.loadLibraries(), nil
}

func (db *DatabaseInMemory) GetLibrariesByTeamIDAndUserID(teamID, userID string) ([]*model.Library, error) {
	allLibs := db.loadLibraries()
	userLibs := make([]*model.Library, 0)
	for _, v := range allLibs {
		if userID == v.OwnerID && teamID == v.TeamID {
			userLibs = append(userLibs, v)
		}
	}
	return userLibs, nil
}

func (db *DatabaseInMemory) AddLibrary(l *model.Library) error {
	l.ID = uuid.NewV4().String()
	allLibs := append(db.loadLibraries(), l)
	db.saveLibraries(allLibs)
	return nil
}

func (db *DatabaseInMemory) SaveLibrary(c *gin.Context, l *model.Library) error {
	allLibs := db.loadLibraries()
	for i, v := range allLibs {
		if l.ID == v.ID {
			allLibs[i] = l
			break
		}
	}
	db.saveLibraries(allLibs)
	return nil
}
