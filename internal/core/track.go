package core

import (
	"net/http"
)

type Track struct {
	ID      string
	Name    string
	Artist  string
	Album   string
	Artwork string
	Genre   string
	Year    string
}

type TrackStorage interface {
	SearchTracks(r *http.Request, filters []Filter) ([]Track, error)
	SaveTracks(r *http.Request, tracks []Track, name, desc string) error
}
