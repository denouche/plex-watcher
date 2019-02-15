package notifications

type Notifications interface {
	SendMessage(string) error
	Answer(url, message string) error
}
