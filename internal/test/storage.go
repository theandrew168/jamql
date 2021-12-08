package test

import (
	"net/http"

	"github.com/theandrew168/jamql/internal/core"
)

type mockStorage struct {
	tracks []core.Track
}

func NewMockStorage(tracks []core.Track) core.Storage {
	s := mockStorage{
		tracks: tracks,
	}
	return &s
}

func (s *mockStorage) SearchTracks(r *http.Request, filters []core.Filter) ([]core.Track, error) {
	// handle no filters as a special case (return no tracks)
	if len(filters) == 0 {
		return []core.Track{}, nil
	}

	// apply each filter
	tracks := s.tracks
	for _, filter := range filters {
		tracks = filter.Apply(tracks)
	}

	return tracks, nil
}

func (s *mockStorage) SaveTracks(r *http.Request, tracks []core.Track, name, desc string) error {
	// mocking save is a no-op
	return nil
}
