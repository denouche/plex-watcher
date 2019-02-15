package plex

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/denouche/plex-watcher/storage/model"

	"github.com/denouche/plex-watcher/utils"
)

const (
	streamTypeVideo    = 1
	streamTypeAudio    = 2
	streamTypeSubtitle = 3
)

type PlexHTTP struct {
	httpClient *http.Client
}

func NewPlex() Plex {
	return &PlexHTTP{
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (s *PlexHTTP) request(url string) (*model.PlexResponse, error) {
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	r.Header.Set(utils.HeaderNameAccept, utils.HeaderValueApplicationJSON)

	resp, err := s.httpClient.Do(r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("response status not OK")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	plexResp := &model.PlexResponse{}
	err = json.Unmarshal(body, plexResp)
	if err != nil {
		return nil, err
	}

	return plexResp, nil

}

func (s *PlexHTTP) GetRecentlyAdded(library *model.Library) (*model.PlexResponse, error) {
	return s.request(fmt.Sprintf("%s/library/recentlyAdded?X-Plex-Token=%s", library.BaseURL, library.Token))
}

func (s *PlexHTTP) GetTitles(library *model.Library, metadata *model.PlexMetadata) []string {
	switch metadata.Type {
	case "movie":
		details := s.getMediaDetails(library, metadata.Key)
		return []string{metadata.Title + "\n" + details}
	case "season":
		datas := s.getMetadataAddedAt(library, metadata.Key, metadata.AddedAt)
		titles := make([]string, 0)
		for _, data := range datas {
			details := s.getMediaDetails(library, data.Key)
			titles = append(titles, fmt.Sprintf("%s - %s - Ep %d \n %s", metadata.ParentTitle, metadata.Title, data.Index, details))
		}
		return titles
	default:
		utils.GetLogger().WithField("type", metadata.Type).Warn("unknown plex metadata type")
	}
	return []string{metadata.Title}
}

func (s *PlexHTTP) getMediaDetails(library *model.Library, url string) string {
	resp, err := s.request(fmt.Sprintf("%s%s?X-Plex-Token=%s", library.BaseURL, url, library.Token))
	if err != nil {
		utils.GetLogger().WithError(err).Warn("error while getting episode index")
		return ""
	}

	if len(resp.MediaContainer.Metadata) > 0 &&
		len(resp.MediaContainer.Metadata[0].Media) > 0 &&
		len(resp.MediaContainer.Metadata[0].Media[0].Part) > 0 {

		media := resp.MediaContainer.Metadata[0].Media[0]
		streams := resp.MediaContainer.Metadata[0].Media[0].Part[0].Stream

		videoTitle := ""
		audios := make([]string, 0)
		subtitles := make([]string, 0)
		for _, stream := range streams {
			switch stream.StreamType {
			case streamTypeVideo:
				videoTitle = stream.DisplayTitle
			case streamTypeAudio:
				audios = append(audios, stream.DisplayTitle)
			case streamTypeSubtitle:
				subtitles = append(subtitles, stream.DisplayTitle)
			}
		}
		details := fmt.Sprintf("%s @%d kbps \nAudio: %s \nSubtitles: %s",
			videoTitle,
			media.Bitrate,
			strings.Join(audios, " / "),
			strings.Join(subtitles, " / "))
		return details
	}
	return ""
}

func (s *PlexHTTP) getMetadataAddedAt(library *model.Library, url string, addedAt *model.Timestamp) []*model.PlexMetadata {
	resp, err := s.request(fmt.Sprintf("%s%s?X-Plex-Token=%s", library.BaseURL, url, library.Token))
	if err != nil {
		utils.GetLogger().WithError(err).Warn("error while getting episode index")
		return make([]*model.PlexMetadata, 0)
	}

	datas := make([]*model.PlexMetadata, 0)
	for _, v := range resp.MediaContainer.Metadata {
		if time.Time(*v.AddedAt).Unix() == time.Time(*addedAt).Unix() {
			datas = append(datas, v)
		}
	}
	return datas
}
