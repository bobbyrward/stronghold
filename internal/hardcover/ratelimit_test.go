package hardcover

import (
	"context"
	"net/http"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

// TestLimiterSpacesAfterBurst verifies the configured limiter paces requests once
// the burst is exhausted: the next reservation must wait roughly the per-request
// interval rather than firing immediately.
func TestLimiterSpacesAfterBurst(t *testing.T) {
	limiter := rate.NewLimiter(rate.Every(time.Minute/requestsPerMinute), burst)

	// Drain the burst — these are immediate.
	for i := 0; i < burst; i++ {
		if d := limiter.Reserve().Delay(); d != 0 {
			t.Fatalf("burst token %d should be immediate, waited %v", i, d)
		}
	}

	// Next one must wait close to the interval (one token / requestsPerMinute).
	interval := time.Minute / requestsPerMinute
	if d := limiter.Reserve().Delay(); d <= 0 || d > interval {
		t.Fatalf("post-burst reservation delay %v, want (0, %v]", d, interval)
	}
}

// TestRateLimitedTransportRespectsContext ensures the transport surfaces a
// cancelled context instead of blocking forever on Wait.
func TestRateLimitedTransportRespectsContext(t *testing.T) {
	// Real burst, then drain it: the next Wait must block on the limiter (one
	// token per hour), so only context cancellation can unblock it. With burst 0
	// Wait would error on "exceeds burst" before ever consulting the context.
	limiter := rate.NewLimiter(rate.Every(time.Hour), 1)
	limiter.Allow() // consume the single burst token
	tr := &rateLimitedTransport{limiter: limiter, base: http.DefaultTransport}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.com", nil)
	if _, err := tr.RoundTrip(req); err == nil {
		t.Fatal("expected error from cancelled context, got nil")
	}
}
