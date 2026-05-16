package pokecache

import (
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	cache := NewCache(time.Minute)
	key := "https://example.com"
	val := []byte("test data")

	cache.Add(key, val)

	got, ok := cache.Get(key)
	if !ok {
		t.Fatal("expected cache hit, got miss")
	}
	if string(got) != string(val) {
		t.Errorf("got %q, want %q", got, val)
	}
}

func TestGetMiss(t *testing.T) {
	cache := NewCache(time.Minute)

	_, ok := cache.Get("https://missing.com")
	if ok {
		t.Fatal("expected cache miss, got hit")
	}
}

func TestReap(t *testing.T) {
	interval := 50 * time.Millisecond
	cache := NewCache(interval)

	cache.Add("https://example.com", []byte("data"))

	// entry should still be present before interval elapses
	if _, ok := cache.Get("https://example.com"); !ok {
		t.Fatal("entry should exist before expiry")
	}

	// wait long enough for two reap ticks
	time.Sleep(interval * 3)

	if _, ok := cache.Get("https://example.com"); ok {
		t.Fatal("entry should have been reaped after expiry")
	}
}
