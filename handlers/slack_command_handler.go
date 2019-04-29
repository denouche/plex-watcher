package handlers

import (
	"net/http"
	"regexp"

	"github.com/denouche/plex-watcher/storage/model"
	"github.com/denouche/plex-watcher/utils"
	"github.com/gin-gonic/gin"
)

const (
	slackParameterTeamID      = "team_id"
	slackParameterUserID      = "user_id"
	slackParameterChannelID   = "channel_id"
	slackParameterText        = "text"
	slackParameterResponseURL = "response_url"
)

var (
	reRegister = regexp.MustCompile(`^register ([^\s]+) ([^\s]+)$`)
	reList     = regexp.MustCompile("^list$")
	reHelp     = regexp.MustCompile("^(?:help|aide)$")
)

func (hc *handlersContext) HandleCommand(c *gin.Context) {
	text, _ := c.GetPostForm(slackParameterText)

	switch {
	case reRegister.MatchString(text):
		utils.GetLoggerFromCtx(c).Debug("Register")
		hc.handleRegister(c)
	case reList.MatchString(text):
		utils.GetLoggerFromCtx(c).Debug("List")
		hc.handleList(c)
	default:
		utils.GetLoggerFromCtx(c).Debug("Unknown command", "command", text)
		responseURL, _ := c.GetPostForm(slackParameterResponseURL)
		_ = hc.notif.Answer(responseURL, "Unknown command")
		fallthrough
	case reHelp.MatchString(text):
		hc.handleHelp(c)
	}

	utils.JSON(c.Writer, http.StatusOK, nil)
}

func (hc *handlersContext) handleRegister(c *gin.Context) {
	teamID, _ := c.GetPostForm(slackParameterTeamID)
	channelID, _ := c.GetPostForm(slackParameterChannelID)
	userID, _ := c.GetPostForm(slackParameterUserID)
	text, _ := c.GetPostForm(slackParameterText)
	responseURL, _ := c.GetPostForm(slackParameterResponseURL)

	utils.GetLoggerFromCtx(c).
		WithField("teamID", teamID).
		WithField("channelID", channelID).
		WithField("userID", userID).
		WithField("message", text).
		WithField("responseURL", responseURL).
		Debug("handleRegister")

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
		_ = hc.notif.Answer(responseURL, "An error occurred while registering library")
	} else {
		_ = hc.notif.Answer(responseURL, "Your library is registered")
	}
}

func (hc *handlersContext) handleList(c *gin.Context) {
	teamID, _ := c.GetPostForm(slackParameterTeamID)
	userID, _ := c.GetPostForm(slackParameterUserID)
	responseURL, _ := c.GetPostForm(slackParameterResponseURL)
	libs, err := hc.db.GetLibrariesByTeamIDAndUserID(teamID, userID)
	if err == nil {
		result := "Your registered libraries are:"
		for _, lib := range libs {
			result = result + "\n" + lib.BaseURL
		}
		_ = hc.notif.Answer(responseURL, result)
	}
}

func (hc *handlersContext) handleHelp(c *gin.Context) {
	responseURL, _ := c.GetPostForm(slackParameterResponseURL)
	_ = hc.notif.Answer(responseURL, "TODO aide")
}
