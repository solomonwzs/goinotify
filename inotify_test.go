package goinotify

import (
	"testing"
	"time"
)

func TestInotify(t *testing.T) {
	w, err := NewWatcher(0)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	w.Close()
	time.Sleep(1 * time.Second)
}
