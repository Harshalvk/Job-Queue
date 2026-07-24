package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/harshalvk/kairos/internal/ratelimit"
)

func TestWait_UnlimitedTypePassesImmediately(t *testing.T) {
	l := ratelimit.New()
	ctx := context.Background()

	start := time.Now()
	err := l.Wait(ctx, "no_limit_configured")
	assert.NoError(t, err)
	assert.Less(t, time.Since(start), 50*time.Millisecond)
}

func TestWait_ThrottlesBeyondBurst(t *testing.T) {
	l := ratelimit.New()
	l.SetLimit("send_email", 2, 1) // 2/sec, burst of 1
	ctx := context.Background()

	// first call should consume the burst token immediately
	start := time.Now()
	assert.NoError(t, l.Wait(ctx, "send_email"))
	assert.Less(t, time.Since(start), 50*time.Millisecond)

	// second call has no burst left, must wait ~500ms for the next token
	start = time.Now()
	assert.NoError(t, l.Wait(ctx, "send_email"))
	assert.GreaterOrEqual(t, time.Since(start), 400*time.Millisecond)
}

func TestWait_RespectsContextCancellation(t *testing.T) {
	l := ratelimit.New()
	l.SetLimit("send_email", 1, 1) // 1/sec, burst 1
	ctx := context.Background()

	require := assert.New(t)
	require.NoError(l.Wait(ctx, "send_email")) // consume the only burst token

	cancelCtx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := l.Wait(cancelCtx, "send_email")
	require.Error(err) // should time out waiting for the next token
}
