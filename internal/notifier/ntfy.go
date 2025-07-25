package notifier

import (
	"net/http"
	"strings"
)

type Noop struct{}

func (n Noop) Notify(_, _ string) error {
	return nil
}

func NewNoop() Noop {
	return Noop{}
}

type Ntfy struct {
	url string
}

func NewNtfy(url string) *Ntfy {
	return &Ntfy{url: url}
}

func (n Ntfy) Notify(title, message string) error {
	req, _ := http.NewRequest("POST", n.url, strings.NewReader(message))
	req.Header.Set("Title", title)
	req.Header.Set("Tags", "warning,notification")
	_, err := http.DefaultClient.Do(req)

	return err
}
