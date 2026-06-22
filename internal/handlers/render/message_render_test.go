package render

import (
	"bbs-go/internal/models"
	"testing"
)

func TestMessageDetailURLFallbackDoesNotPointToRemovedMessagesPage(t *testing.T) {
	got := getMessageDetailUrl(&models.Message{Type: 999})
	if got != "/user/profile" {
		t.Fatalf("expected unknown message fallback to profile, got %q", got)
	}
}
