package model

import (
	"fmt"
	"strconv"
	"time"
)

type PlexResponse struct {
	MediaContainer PlexMediaContainer `json:"MediaContainer"`
}

type PlexMediaContainer struct {
	Metadata []*PlexMetadata `json:"Metadata"`
}

type PlexMetadata struct {
	ParentTitle string       `json:"parentTitle"`
	Title       string       `json:"title"`
	Index       int          `json:"index"`
	Type        string       `json:"type"`
	AddedAt     *Timestamp   `json:"addedAt"`
	Key         string       `json:"key"`
	Media       []*PlexMedia `json:"Media"`
}

type PlexMedia struct {
	VideoResolution string           `json:"videoResolution"`
	Bitrate         int              `json:"bitrate"`
	Width           int              `json:"width"`
	Height          int              `json:"height"`
	VideoCodec      string           `json:"videoCodec"`
	Part            []*PlexMediaPart `json:"Part"`
}

type PlexMediaPart struct {
	Stream []*PlexMediaPartStream `json:"Stream"`
}

type PlexMediaPartStream struct {
	StreamType   int    `json:"streamType"`
	DisplayTitle string `json:"displayTitle"'`
}

type Timestamp time.Time

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	ts := time.Time(*t).Unix()
	stamp := fmt.Sprint(ts)
	return []byte(stamp), nil
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	ts, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}
	*t = Timestamp(time.Unix(int64(ts), 0))
	return nil
}
