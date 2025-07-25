package ntfy

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/egregors/hk/log"
)

// Config holds configuration for ntfy notifications
type Config struct {
	URL     string // The ntfy.sh URL with topic
	Enabled bool   // Whether notifications are enabled
}

// Notifier handles sending notifications to ntfy.sh
type Notifier struct {
	config Config
	client *http.Client
}

// NoopNotifier is a no-op implementation for when notifications are disabled
type NoopNotifier struct{}

// SendError does nothing for the no-op notifier
func (n *NoopNotifier) SendError(message string) error {
	return nil
}

// New creates a new Notifier with the given configuration
func New(config Config) *Notifier {
	return &Notifier{
		config: config,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendError sends an error notification to ntfy.sh
func (n *Notifier) SendError(message string) error {
	if !n.config.Enabled || n.config.URL == "" {
		return nil
	}

	title := "🇭🇰 Sensor Error"
	body := fmt.Sprintf("Sensor error occurred: %s", message)

	req, err := http.NewRequest("POST", n.config.URL, bytes.NewBufferString(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Title", title)
	req.Header.Set("Priority", "high")
	req.Header.Set("Tags", "warning,sensor")

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("ntfy request failed with status: %d", resp.StatusCode)
	}

	log.Info.Printf("sent ntfy notification: %s", message)
	return nil
}