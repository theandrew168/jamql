package core

import (
	"testing"
)

func TestValidFilters(t *testing.T) {
	tests := []struct {
		Name   string
		Filter Filter
		Count  int
	}{
		{
			Name:   "NameEquals",
			Filter: Filter{"name", "equals", "teflon"},
			Count:  1,
		},
		{
			Name:   "NameContains",
			Filter: Filter{"name", "contains", "land"},
			Count:  2,
		},
		{
			Name:   "ArtistEquals",
			Filter: Filter{"artist", "equals", "ben folds"},
			Count:  4,
		},
		{
			Name:   "ArtistContains",
			Filter: Filter{"artist", "contains", "folds"},
			Count:  7,
		},
		{
			Name:   "AlbumEquals",
			Filter: Filter{"album", "equals", "octahedron"},
			Count:  3,
		},
		{
			Name:   "AlbumContains",
			Filter: Filter{"album", "contains", "amen"},
			Count:  3,
		},
		{
			Name:   "GenreEquals",
			Filter: Filter{"genre", "equals", "prog rock"},
			Count:  3,
		},
		{
			Name:   "GenreContains",
			Filter: Filter{"genre", "contains", "rock"},
			Count:  10,
		},
		{
			Name:   "YearEquals",
			Filter: Filter{"year", "equals", "1997"},
			Count:  3,
		},
		{
			Name:   "YearBetween",
			Filter: Filter{"year", "between", "1990-1999"},
			Count:  3,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			matches := test.Filter.Apply(sampleTracks)
			count := len(matches)
			if count != test.Count {
				t.Errorf("want %d; got %d", test.Count, count)
			}
		})
	}
}
