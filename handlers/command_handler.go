package handlers

import (
	"net/http"
	"regexp"

	"github.com/denouche/plex-watcher/storage/model"
	"github.com/denouche/plex-watcher/utils"
	"github.com/gin-gonic/gin"
)

var (
	reRegister = regexp.MustCompile(`^register ([^\s]+) ([^\s]+)$`)
	reList     = regexp.MustCompile("^list$")
)

func (hc *handlersContext) Command(c *gin.Context) {
	userID, _ := c.GetPostForm("user_id")
	channelID, _ := c.GetPostForm("channel_id")
	teamID, _ := c.GetPostForm("team_id")
	text, _ := c.GetPostForm("text")
	responseURL, _ := c.GetPostForm("response_url")

	utils.GetLoggerFromCtx(c).
		WithField("teamID", teamID).
		WithField("channelID", channelID).
		WithField("userID", userID).
		WithField("message", text).
		WithField("responseURL", responseURL).
		Debug("command")

	switch {
	case reRegister.MatchString(text):
		utils.GetLoggerFromCtx(c).Debug("Register")
		matcher := reRegister.FindStringSubmatch(text)
		lib := &model.Library{
			OwnerID:   userID,
			TeamID:    teamID,
			ChannelID: channelID,
			BaseURL:   matcher[1],
			Token:     matcher[2],
		}
		err := hc.db.AddLibrary(lib)
		if err != nil {
			utils.GetLoggerFromCtx(c).Error("error occured while registering lib", "error", err)
			hc.notif.Answer(responseURL, "An error occurred while registering library")
		} else {
			hc.notif.Answer(responseURL, "Your library is registered")
		}
	case reList.MatchString(text):
		utils.GetLoggerFromCtx(c).Debug("List")
		libs, err := hc.db.GetLibrariesByTeamIDAndUserID(teamID, userID)
		if err == nil {
			result := "Your registered libraries are:"
			for _, lib := range libs {
				result = result + "\n" + lib.BaseURL
			}
			_ = hc.notif.Answer(responseURL, result)
		}
	default:
		utils.GetLoggerFromCtx(c).Debug("Unknown command", "command", text)
		_ = hc.notif.Answer(responseURL, "Unknown command")
	}

	utils.JSON(c.Writer, http.StatusOK, nil)
}
