package core

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
	SearchTracks(filters []Filter) ([]Track, error)
	SaveTracks(tracks []Track, name, desc string) error
}
