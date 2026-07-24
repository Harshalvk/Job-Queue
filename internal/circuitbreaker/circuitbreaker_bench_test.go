package circuitbreaker_test

import (
	"testing"
	"time"

	"github.com/harshalvk/kairos/internal/circuitbreaker"
)

func BenchmarkAllow_Closed(b *testing.B) {
	cb := circuitbreaker.New(5, time.Second)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cb.Allow("send_email")
	}
}

func BenchmarkRecordFailure(b *testing.B) {
	cb := circuitbreaker.New(1000000, time.Second) // high threshold, never trips mid-benchmark

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cb.RecordFailure("send_email")
	}
}

func BenchmarkAllow_Parallel(b *testing.B) {
	cb := circuitbreaker.New(5, time.Second)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cb.Allow("send_email")
		}
	})
}
