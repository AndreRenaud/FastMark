package main

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/AndreRenaud/fastmark/storage"
)

func writeTestImage(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, image.NewRGBA(image.Rect(0, 0, 2, 2))); err != nil {
		t.Fatal(err)
	}
}

func waitReady(t *testing.T, c *imageCache, name string) *imageEntry {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		c.mu.Lock()
		entry := c.entries[name]
		c.mu.Unlock()
		if entry != nil && entry.ready(false) {
			return entry
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatalf("%s was not read ahead", name)
	return nil
}

func TestImageCacheReadAhead(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"a.png", "b.png", "c.png"} {
		writeTestImage(t, filepath.Join(dir, "images", name))
	}
	c := newImageCache(storage.NewStorage(dir))

	// The first entry is left for the caller to load; the rest are read ahead.
	first := c.retain([]string{"a.png", "b.png", "missing.png"}, false)
	if first == nil || first.ready(false) {
		t.Fatalf("first entry should be returned unloaded, got %+v", first)
	}
	first.load()
	if first.err != nil || first.img == nil {
		t.Fatalf("load a.png: img=%v err=%v", first.img, first.err)
	}
	b := waitReady(t, c, "b.png")
	if b.err != nil {
		t.Fatalf("read ahead b.png: %s", b.err)
	}
	if missing := waitReady(t, c, "missing.png"); missing.err == nil {
		t.Fatal("missing.png should have failed to load")
	}

	// Moving on reuses what was read ahead, drops what left the window and
	// retries failures.
	next := c.retain([]string{"b.png", "c.png", "missing.png"}, false)
	if next != b {
		t.Fatal("b.png was not reused from the cache")
	}
	c.mu.Lock()
	_, kept := c.entries["a.png"]
	retried := c.entries["missing.png"]
	c.mu.Unlock()
	if kept {
		t.Fatal("a.png should have been dropped from the cache")
	}
	if retried == nil || retried.ready(false) && retried.err == nil {
		t.Fatal("missing.png should have been retried")
	}
	waitReady(t, c, "c.png")
}

func TestReadAheadWindow(t *testing.T) {
	m := &appModel{files: []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}}

	m.selectedIndex, m.navDirection = 5, 1
	if got, want := m.readAheadWindow(), []string{"5", "6", "4", "7", "3", "8"}; !equalStrings(got, want) {
		t.Errorf("forwards: got %v, want %v", got, want)
	}
	m.selectedIndex, m.navDirection = 5, -1
	if got, want := m.readAheadWindow(), []string{"5", "4", "6", "3", "7", "2"}; !equalStrings(got, want) {
		t.Errorf("backwards: got %v, want %v", got, want)
	}
	m.selectedIndex, m.navDirection = 0, 0
	if got, want := m.readAheadWindow(), []string{"0", "1", "2", "3"}; !equalStrings(got, want) {
		t.Errorf("start: got %v, want %v", got, want)
	}
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
