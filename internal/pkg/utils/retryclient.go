package utils

import (
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/go-retryablehttp"
)

func NewExponentialBackoff() retryablehttp.Backoff {
	// Create an exponential backoff policy with jitter.
	// These values can be tuned to your needs.
	b := &backoff.ExponentialBackOff{
		InitialInterval:     100 * time.Millisecond,
		RandomizationFactor: 0.5, // Adds jitter
		Multiplier:          1.5,
		MaxInterval:         1 * time.Second,
		MaxElapsedTime:      4 * time.Second,
		Clock:               backoff.SystemClock,
	}
	b.Reset() // Important to start the backoff fresh

	// The returned function is a closure that captures the backoff policy.
	return func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		// The cenkalti/backoff library handles the backoff duration calculation.
		// It internally tracks the number of retries.
		// The min, max, and attemptNum from retryablehttp are ignored in this case,
		// as the cenkalti/backoff policy manages its own state.
		return b.NextBackOff()
	}
}
