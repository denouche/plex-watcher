package notifications

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/denouche/plex-watcher/utils"
)

type Slack struct {
	url        string
	httpClient *http.Client
}

func NewSlack(url string) Notifications {
	return &Slack{
		url: url,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (s *Slack) Answer(url, message string) error {
	message = strings.Replace(message, `\`, `\\`, -1)
	message = strings.Replace(message, `"`, `\"`, -1)

	messageReader := strings.NewReader(fmt.Sprintf(`{"text":"%s"}`, message))
	r, err := http.NewRequest(http.MethodPost, url, messageReader)
	if err != nil {
		return err
	}

	r.Header.Set(utils.HeaderNameContentType, utils.HeaderValueApplicationJSON)

	resp, err := s.httpClient.Do(r)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("response status not OK")
	}

	return nil
}

func (s *Slack) SendMessage(message string) error {
	return s.Answer(s.url, message)
}
