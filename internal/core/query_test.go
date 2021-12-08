package core

import (
	"testing"
)

func TestQueryBuild(t *testing.T) {
	tests := []struct {
		name  string
		query Query
		want  string
	}{
		{
			name:  "Name",
			query: Query{Name: "Late"},
			want:  `track:"Late"`,
		},
		{
			name:  "Artist",
			query: Query{Artist: "Ben Folds"},
			want:  `artist:"Ben Folds"`,
		},
		{
			name:  "Album",
			query: Query{Album: "Songs for Silverman"},
			want:  `album:"Songs for Silverman"`,
		},
		{
			name:  "Genre",
			query: Query{Genre: "country"},
			want:  `genre:"country"`,
		},
		{
			name:  "Year",
			query: Query{Year: "2000"},
			want:  `year:2000`,
		},
		{
			name:  "YearRange",
			query: Query{Year: "2000-2010"},
			want:  `year:2000-2010`,
		},
		{
			name: "Multiple",
			query: Query{
				Artist: "Ben Folds",
				Album:  "Songs for Silverman",
			},
			want: `artist:"Ben Folds" album:"Songs for Silverman"`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			q := test.query.Build()
			if q != test.want {
				t.Errorf("want %q; got %q", test.want, q)
			}
		})
	}
}

func TestNewQuery(t *testing.T) {
	tests := []struct {
		name    string
		filters []Filter
		want    string
	}{
		{
			name:    "Empty",
			filters: []Filter{},
			want:    ``,
		},
		{
			name: "Single",
			filters: []Filter{
				{"artist", "equals", "Ben Folds"},
			},
			want: `artist:"Ben Folds"`,
		},
		{
			name: "Duplicate",
			filters: []Filter{
				{"artist", "equals", "Ben Folds"},
				{"artist", "equals", "The Mars Volta"},
			},
			want: `artist:"The Mars Volta"`,
		},
		{
			name: "IgnoreInvalidYear",
			filters: []Filter{
				{"artist", "equals", "Ben Folds"},
				{"year", "equals", "InvalidYear"},
			},
			want: `artist:"Ben Folds"`,
		},
		{
			name: "CountryBangers",
			filters: []Filter{
				{"genre", "equals", "country"},
				{"year", "equals", "1990-1999"},
			},
			want: `genre:"country" year:1990-1999`,
		},
		{
			name: "All",
			filters: []Filter{
				{"name", "equals", "Late"},
				{"artist", "equals", "Ben Folds"},
				{"album", "equals", "Songs for Silverman"},
				{"genre", "equals", "Alt Rock"},
				{"year", "equals", "2005"},
			},
			want: `track:"Late" artist:"Ben Folds" album:"Songs for Silverman" genre:"Alt Rock" year:2005`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			q := NewQuery(test.filters).Build()
			if q != test.want {
				t.Errorf("want %q; got %q", test.want, q)
			}
		})
	}
}
