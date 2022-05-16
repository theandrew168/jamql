package test

import (
	"net/http"

	"github.com/theandrew168/jamql/internal/core"
)

type storage struct {
	tracks []core.Track
}

func NewTrackStorage() core.TrackStorage {
	s := storage{
		tracks: SampleTracks,
	}
	return &s
}

func (s *storage) SearchTracks(r *http.Request, filters []core.Filter) ([]core.Track, error) {
	// handle no filters as a special case (return no tracks)
	if len(filters) == 0 {
		return nil, nil
	}

	// apply each filter
	tracks := s.tracks
	for _, filter := range filters {
		tracks = filter.Apply(tracks)
	}

	return tracks, nil
}

func (s *storage) SaveTracks(r *http.Request, tracks []core.Track, name, desc string) error {
	// mocking save is a no-op
	return nil
}
