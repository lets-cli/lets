package lsp

import (
	"sync"
	"time"
)

type documentDebouncer struct {
	delay   time.Duration
	refresh func(string)

	mu     sync.Mutex
	timers map[string]*time.Timer
}

func newDocumentDebouncer(delay time.Duration, refresh func(string)) *documentDebouncer {
	return &documentDebouncer{
		delay:   delay,
		refresh: refresh,
		timers:  make(map[string]*time.Timer),
	}
}

func (d *documentDebouncer) Schedule(uri string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if timer, ok := d.timers[uri]; ok {
		timer.Stop()
	}

	var timer *time.Timer

	timer = time.AfterFunc(d.delay, func() {
		d.fire(uri, timer)
	})

	d.timers[uri] = timer
}

func (d *documentDebouncer) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()

	for uri, timer := range d.timers {
		timer.Stop()
		delete(d.timers, uri)
	}
}

func (d *documentDebouncer) fire(uri string, timer *time.Timer) {
	d.mu.Lock()

	current, ok := d.timers[uri]
	if !ok || current != timer {
		d.mu.Unlock()
		return
	}

	delete(d.timers, uri)
	d.mu.Unlock()

	d.refresh(uri)
}
