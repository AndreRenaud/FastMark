package main

import (
	"image"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/AndreRenaud/fastmark/storage"
)

const (
	// readAheadCount and readBehindCount are how many images either side of
	// the selection (relative to the direction of travel) are kept decoded.
	// Decoded images are large, so these bound memory use.
	readAheadCount  = 3
	readBehindCount = 2

	// readAheadWorkers limits concurrent read-ahead loads so they don't
	// starve the image actually being displayed of bandwidth or CPU.
	readAheadWorkers = 2
)

// imageEntry is one image in the cache, loaded at most once.
type imageEntry struct {
	fetch func() (image.Image, error)

	loadOnce sync.Once
	loaded   atomic.Bool
	img      image.Image // valid once loaded
	err      error       // valid once loaded

	contrastOnce   sync.Once
	contrastLoaded atomic.Bool
	contrast       image.Image
}

// load fetches and decodes the image, blocking until it is available. It is
// safe to call concurrently; only the first call does the work.
func (e *imageEntry) load() {
	e.loadOnce.Do(func() {
		e.img, e.err = e.fetch()
		e.loaded.Store(true)
	})
}

// display returns the image to show, computing (and caching) the
// contrast-stretched version when needed. It must only be called after a
// successful load.
func (e *imageEntry) display(autoContrast bool) image.Image {
	if !autoContrast {
		return e.img
	}
	e.contrastOnce.Do(func() {
		e.contrast = autoContrastImage(e.img)
		e.contrastLoaded.Store(true)
	})
	return e.contrast
}

// ready reports whether load and display(autoContrast) will return without
// blocking.
func (e *imageEntry) ready(autoContrast bool) bool {
	if !e.loaded.Load() {
		return false
	}
	return e.err != nil || !autoContrast || e.contrastLoaded.Load()
}

// imageCache keeps decoded images for a small window of files around the
// current selection, reading ahead so that moving to a neighbouring image
// doesn't have to wait on storage or decoding.
type imageCache struct {
	backend storage.Storage

	mu      sync.Mutex
	entries map[string]*imageEntry
	workers chan struct{}
}

func newImageCache(backend storage.Storage) *imageCache {
	return &imageCache{
		backend: backend,
		entries: make(map[string]*imageEntry),
		workers: make(chan struct{}, readAheadWorkers),
	}
}

// retain makes window (filenames relative to images/, most wanted first) the
// set of cached images: everything else is dropped, and anything missing is
// read ahead in the background. It returns the entry for window[0], which
// the caller is expected to load itself so that it isn't queued behind the
// read-ahead; nil if window is empty.
func (c *imageCache) retain(window []string, autoContrast bool) *imageEntry {
	c.mu.Lock()
	defer c.mu.Unlock()

	old := c.entries
	c.entries = make(map[string]*imageEntry, len(window))
	var first *imageEntry
	for i, name := range window {
		entry, ok := old[name]
		if ok && entry.loaded.Load() && entry.err != nil {
			ok = false // retry failed loads
		}
		if !ok {
			entry = &imageEntry{fetch: func() (image.Image, error) {
				return loadImage(c.backend, filepath.Join("images", name))
			}}
			if i > 0 {
				go c.readAhead(name, entry, autoContrast)
			}
		}
		c.entries[name] = entry
		if i == 0 {
			first = entry
		}
	}
	return first
}

// readAhead loads entry in the background, unless it has dropped out of the
// window by the time a worker slot is free (e.g. when skipping quickly
// through the list).
func (c *imageCache) readAhead(name string, entry *imageEntry, autoContrast bool) {
	c.workers <- struct{}{}
	defer func() { <-c.workers }()

	c.mu.Lock()
	wanted := c.entries[name] == entry
	c.mu.Unlock()
	if !wanted {
		return
	}

	// Warm the region cache too; a missing label file is not an error here.
	LoadRegionList(c.backend, labelFileFor(name))

	entry.load()
	if entry.err == nil {
		entry.display(autoContrast)
	}
}
