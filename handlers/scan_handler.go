package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/denouche/plex-watcher/storage/model"
	"github.com/denouche/plex-watcher/utils"
	"github.com/gin-gonic/gin"
)

func (hc *handlersContext) ScanLibraries(c *gin.Context) {
	allLibs, err := hc.db.GetAllLibraries()
	if err != nil {
		utils.JSON(c.Writer, http.StatusInternalServerError, nil)
		return
	}

	for _, lib := range allLibs {
		utils.GetLoggerFromCtx(c).
			WithField("baseurl", lib.BaseURL).
			WithField("token", lib.Token).
			Debug("scanning lib")
		recentlyAdded, err := hc.plex.GetRecentlyAdded(lib)
		if err != nil {
			utils.GetLoggerFromCtx(c).WithError(err).Error("error while getting recently added")
		}

		if len(recentlyAdded.MediaContainer.Metadata) > 0 {
			if lib.LastNotifiedMediaAddedAt == nil {
				utils.GetLoggerFromCtx(c).Debug("lib never scanned")
				t := time.Time(*recentlyAdded.MediaContainer.Metadata[0].AddedAt)
				lib.LastNotifiedMediaAddedAt = &t
				err = hc.db.SaveLibrary(c, lib)
				if err != nil {
					utils.GetLoggerFromCtx(c).WithError(err).Error("error while saving lib with last nil")
				}
			} else {
				utils.GetLoggerFromCtx(c).WithField("at", lib.LastNotifiedMediaAddedAt).Debug("last lib scan")
				toNotify := make([]*model.PlexMetadata, 0)
				for _, media := range recentlyAdded.MediaContainer.Metadata {
					t := time.Time(*media.AddedAt)
					if lib.LastNotifiedMediaAddedAt.Before(t) {
						toNotify = append(toNotify, media)
					}
				}

				for _, media := range toNotify {
					t := time.Time(*media.AddedAt)
					if lib.LastNotifiedMediaAddedAt.Before(t) {
						lib.LastNotifiedMediaAddedAt = &t
					}
					titles := hc.plex.GetTitles(lib, media)
					for _, title := range titles {
						_ = hc.notif.SendMessage(fmt.Sprintf("New media added to <@%s> library: %s", lib.OwnerID, title))
					}
				}

				err = hc.db.SaveLibrary(c, lib)
				if err != nil {
					utils.GetLoggerFromCtx(c).WithError(err).Error("error while saving lib with last not nil")
				}
			}

		}
	}

	utils.JSON(c.Writer, http.StatusOK, nil)
}
