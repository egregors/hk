package srv

import (
	"testing"
	"time"
)

func TestFormatUptime(t *testing.T) {
	// Create a server with a known start time
	server := &Server{
		startTime: time.Now().Add(-65 * time.Minute), // 1 hour and 5 minutes ago
	}

	uptime := server.formatUptime()
	expected := "(uptime: 1h 5m)"
	
	if uptime != expected {
		t.Errorf("Expected uptime %s, got %s", expected, uptime)
	}
}

func TestFormatUptimeMinutes(t *testing.T) {
	// Test uptime less than an hour
	server := &Server{
		startTime: time.Now().Add(-30 * time.Minute), // 30 minutes ago
	}

	uptime := server.formatUptime()
	expected := "(uptime: 30m)"
	
	if uptime != expected {
		t.Errorf("Expected uptime %s, got %s", expected, uptime)
	}
}

func TestFormatUptimeDays(t *testing.T) {
	// Test uptime more than a day
	server := &Server{
		startTime: time.Now().Add(-25*time.Hour - 30*time.Minute), // 1 day, 1 hour, 30 minutes ago
	}

	uptime := server.formatUptime()
	expected := "(uptime: 1d 1h 30m)"
	
	if uptime != expected {
		t.Errorf("Expected uptime %s, got %s", expected, uptime)
	}
}

func TestTitleWithUptime(t *testing.T) {
	// Test that title includes uptime
	server := &Server{
		sensorStatus: ONLINE,
		startTime:    time.Now().Add(-45 * time.Minute), // 45 minutes ago
	}

	title := server.title()
	expected := "Sensor: 🟢 Online (uptime: 45m)\n"
	
	if title != expected {
		t.Errorf("Expected title %q, got %q", expected, title)
	}
}

func TestTitleOfflineWithUptime(t *testing.T) {
	// Test that title includes uptime even when offline
	server := &Server{
		sensorStatus: OFFLINE,
		sensorErr:    &testError{msg: "test error"},
		startTime:    time.Now().Add(-2*time.Hour - 15*time.Minute), // 2 hours, 15 minutes ago
	}

	title := server.title()
	expected := "Sensor: 🔴 Offline (uptime: 2h 15m)\nError: test error\n"
	
	if title != expected {
		t.Errorf("Expected title %q, got %q", expected, title)
	}
}

// Helper type for testing errors
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}