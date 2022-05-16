package core_test

import (
	"testing"

	"github.com/theandrew168/jamql/pkg/core"
	"github.com/theandrew168/jamql/pkg/test"
)

func TestValidFilters(t *testing.T) {
	tests := []struct {
		name   string
		filter core.Filter
		count  int
	}{
		{
			"NameEquals",
			core.Filter{"name", "equals", "teflon"},
			1,
		},
		{
			"NameContains",
			core.Filter{"name", "contains", "land"},
			2,
		},
		{
			"ArtistEquals",
			core.Filter{"artist", "equals", "ben folds"},
			4,
		},
		{
			"ArtistContains",
			core.Filter{"artist", "contains", "folds"},
			7,
		},
		{
			"AlbumEquals",
			core.Filter{"album", "equals", "octahedron"},
			3,
		},
		{
			"AlbumContains",
			core.Filter{"album", "contains", "amen"},
			3,
		},
		{
			"GenreEquals",
			core.Filter{"genre", "equals", "prog rock"},
			3,
		},
		{
			"GenreContains",
			core.Filter{"genre", "contains", "rock"},
			10,
		},
		{
			"YearEquals",
			core.Filter{"year", "equals", "1997"},
			3,
		},
		{
			"YearBetween",
			core.Filter{"year", "between", "1990-1999"},
			3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := tt.filter.Apply(test.SampleTracks)
			count := len(matches)
			if count != tt.count {
				t.Errorf("want %d; got %d", tt.count, count)
			}
		})
	}
}
