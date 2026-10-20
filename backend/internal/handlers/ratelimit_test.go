package handlers

import "testing"

func TestAuthRateLimiter(t *testing.T) {
	l := newAuthRateLimiter()
	for i := 0; i < authMaxAttempts; i++ {
		if !l.Allow("1.2.3.4") {
			t.Fatalf("attempt %d unexpectedly denied", i+1)
		}
	}
	if l.Allow("1.2.3.4") {
		t.Fatal("attempt over limit should be denied")
	}
	if !l.Allow("5.6.7.8") {
		t.Fatal("other key should not be affected")
	}
}
