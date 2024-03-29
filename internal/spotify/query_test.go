package spotify_test

import (
	"testing"

	"github.com/theandrew168/jamql/internal/core"
	"github.com/theandrew168/jamql/internal/spotify"
)

func TestQueryBuild(t *testing.T) {
	tests := []struct {
		name  string
		query spotify.Query
		want  string
	}{
		{
			"Name",
			spotify.Query{Name: "Late"},
			`track:"Late"`,
		},
		{
			"Artist",
			spotify.Query{Artist: "Ben Folds"},
			`artist:"Ben Folds"`,
		},
		{
			"Album",
			spotify.Query{Album: "Songs for Silverman"},
			`album:"Songs for Silverman"`,
		},
		{
			"Genre",
			spotify.Query{Genre: "country"},
			`genre:"country"`,
		},
		{
			"Year",
			spotify.Query{Year: "2000"},
			`year:2000`,
		},
		{
			"YearRange",
			spotify.Query{Year: "2000-2010"},
			`year:2000-2010`,
		},
		{
			"Multiple",
			spotify.Query{
				Artist: "Ben Folds",
				Album:  "Songs for Silverman",
			},
			`artist:"Ben Folds" album:"Songs for Silverman"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := tt.query.String()
			if q != tt.want {
				t.Errorf("want %q; got %q", tt.want, q)
			}
		})
	}
}

func TestNewQuery(t *testing.T) {
	tests := []struct {
		name    string
		filters []core.Filter
		want    string
	}{
		{
			"Empty",
			[]core.Filter{},
			``,
		},
		{
			"Single",
			[]core.Filter{
				{"artist", "equals", "Ben Folds"},
			},
			`artist:"Ben Folds"`,
		},
		{
			"Duplicate",
			[]core.Filter{
				{"artist", "equals", "Ben Folds"},
				{"artist", "equals", "The Mars Volta"},
			},
			`artist:"The Mars Volta"`,
		},
		{
			"IgnoreInvalidYear",
			[]core.Filter{
				{"artist", "equals", "Ben Folds"},
				{"year", "equals", "InvalidYear"},
			},
			`artist:"Ben Folds"`,
		},
		{
			"CountryBangers",
			[]core.Filter{
				{"genre", "equals", "country"},
				{"year", "equals", "1990-1999"},
			},
			`genre:"country" year:1990-1999`,
		},
		{
			"All",
			[]core.Filter{
				{"name", "equals", "Late"},
				{"artist", "equals", "Ben Folds"},
				{"album", "equals", "Songs for Silverman"},
				{"genre", "equals", "Alt Rock"},
				{"year", "equals", "2005"},
			},
			`track:"Late" artist:"Ben Folds" album:"Songs for Silverman" genre:"Alt Rock" year:2005`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := spotify.NewQuery(tt.filters).String()
			if q != tt.want {
				t.Errorf("want %q; got %q", tt.want, q)
			}
		})
	}
}
