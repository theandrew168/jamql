package spotify

import (
	"net/http"

	"github.com/golangcollege/sessions"

	"github.com/theandrew168/jamql/internal/core"
)

type storage struct {
	session *sessions.Session
}

func NewTrackStorage(session *sessions.Session) core.TrackStorage {
	s := storage{
		session: session,
	}
	return &s
}

func (s *storage) SearchTracks(r *http.Request, filters []core.Filter) ([]core.Track, error) {
	// handle no filters as a special case (return no tracks)
	if len(filters) == 0 {
		return []core.Track{}, nil
	}

//	// apply each filter
//	tracks := s.tracks
//	for _, filter := range filters {
//		tracks = filter.Apply(tracks)
//	}

	return []core.Track{}, nil
}

func (s *storage) SaveTracks(r *http.Request, tracks []core.Track, name, desc string) error {
	return nil
}
