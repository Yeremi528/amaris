// Package healthgrp holds a handler for health checking.
package healthgrp

import (
	"context"
	"dragonball/foundation/timecl"
	"dragonball/foundation/web"
	"net/http"
	"time"
)

const (
	UP = "UP"
)

// Handlers manages the health endpoint used by k8s.
type Handlers struct {
	build string
	since time.Time
	cores string
}

// New constructs a Handlers type for route access.
func New(build string, since time.Time, cores string) *Handlers {
	return &Handlers{
		build: build,
		since: since,
		cores: cores,
	}
}

// Health allows k8s to know if the service is running.
func (h *Handlers) Health(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	now := timecl.Now()
	uptime := now.Sub(h.since).Truncate(time.Second)
	uptimestr := uptime.String()

	data := struct {
		Status     string    `json:"status,omitempty"`
		Uptime     string    `json:"uptime,omitempty"`
		Timestamp  time.Time `json:"timestamp,omitempty"`
		Since      time.Time `json:"since,omitempty"`
		Version    string    `json:"version,omitempty"`
		GOMAXPROCS string    `json:"GOMAXPROCS,omitempty"`
	}{
		Status:     UP,
		Uptime:     uptimestr,
		Timestamp:  now,
		Since:      h.since,
		Version:    h.build,
		GOMAXPROCS: h.cores,
	}

	return web.Respond(ctx, w, data, http.StatusOK)
}
