package spotify

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/theandrew168/jamql/internal/core"
)

var (
	yearPattern = regexp.MustCompile(`^\d{4}(-\d{4})?$`)
)

// Represents the query options available to Spotify's Search API.
//
// https://developer.spotify.com/console/get-search-item/
type Query struct {
	Name   string
	Artist string
	Album  string
	Genre  string
	Year   string
}

// Converts a slice of core.Filter into a Query.
//
// The current behavior is to use the last value provided for a given field.
// This means that if multiple filters with FilterKey == Artist were offered,
// then the last artist value would be sent to Spotify.
//
// This function treats Equal and Contain the same because further
// filtering is expected to be performed somewhere down the line.
func NewQuery(filters []core.Filter) Query {
	var q Query
	for _, filter := range filters {
		switch filter.Key {
		case "name":
			q.Name = filter.Value
		case "artist":
			q.Artist = filter.Value
		case "album":
			q.Album = filter.Value
		case "genre":
			q.Genre = filter.Value
		case "year":
			q.Year = filter.Value
		}
	}

	return q
}

// Assemble the Query's fields into a Spotify-compatible string. If some fields
// aren't present (aka they are their zero value) then they won't be part of
// the built query string.
func (q Query) String() string {
	query := []string{}
	if q.Name != "" {
		query = append(query, fmt.Sprintf(`track:"%s"`, q.Name))
	}
	if q.Artist != "" {
		query = append(query, fmt.Sprintf(`artist:"%s"`, q.Artist))
	}
	if q.Album != "" {
		query = append(query, fmt.Sprintf(`album:"%s"`, q.Album))
	}
	if q.Genre != "" {
		query = append(query, fmt.Sprintf(`genre:"%s"`, q.Genre))
	}
	if q.Year != "" {
		if yearPattern.MatchString(q.Year) {
			query = append(query, fmt.Sprintf(`year:%s`, q.Year))
		}
	}

	return strings.Join(query, " ")
}
