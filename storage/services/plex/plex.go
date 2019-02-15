package plex

import "github.com/denouche/plex-watcher/storage/model"

type Plex interface {
	GetRecentlyAdded(library *model.Library) (*model.PlexResponse, error)
	GetTitles(library *model.Library, metadata *model.PlexMetadata) []string
}
