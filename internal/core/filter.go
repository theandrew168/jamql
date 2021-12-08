package core

import (
	"strconv"
	"strings"
)

type Filter struct {
	Key   string
	Op    string
	Value string
}

// Name
func filterNameEquals(track Track, value string) bool {
	return strings.ToLower(track.Name) == strings.ToLower(value)
}

func filterNameContains(track Track, value string) bool {
	return strings.Contains(strings.ToLower(track.Name), strings.ToLower(value))
}

// Artist
func filterArtistEquals(track Track, value string) bool {
	return strings.ToLower(track.Artist) == strings.ToLower(value)
}

func filterArtistContains(track Track, value string) bool {
	return strings.Contains(strings.ToLower(track.Artist), strings.ToLower(value))
}

// Album
func filterAlbumEquals(track Track, value string) bool {
	return strings.ToLower(track.Album) == strings.ToLower(value)
}

func filterAlbumContains(track Track, value string) bool {
	return strings.Contains(strings.ToLower(track.Album), strings.ToLower(value))
}

// Genre
func filterGenreEquals(track Track, value string) bool {
	return strings.ToLower(track.Genre) == strings.ToLower(value)
}

func filterGenreContains(track Track, value string) bool {
	return strings.Contains(strings.ToLower(track.Genre), strings.ToLower(value))
}

// Year
func filterYearEquals(track Track, value string) bool {
	return strings.ToLower(track.Year) == strings.ToLower(value)
}

func filterYearBetween(track Track, value string) bool {
	years := strings.Split(value, "-")
	start, err := strconv.Atoi(years[0])
	if err != nil {
		return false
	}
	end, err := strconv.Atoi(years[1])
	if err != nil {
		return false
	}

	year, err := strconv.Atoi(track.Year)
	if err != nil {
		return false
	}

	return year >= start && year <= end
}

// matrix of specific filter functions
type filterFunc func(Track, string) bool

var filterFuncs = map[string]map[string]filterFunc{
	"name": {
		"equals":   filterNameEquals,
		"contains": filterNameContains,
	},
	"artist": {
		"equals":   filterArtistEquals,
		"contains": filterArtistContains,
	},
	"album": {
		"equals":   filterAlbumEquals,
		"contains": filterAlbumContains,
	},
	"genre": {
		"equals":   filterGenreEquals,
		"contains": filterGenreContains,
	},
	"year": {
		"equals":  filterYearEquals,
		"between": filterYearBetween,
	},
}

func (f Filter) Apply(tracks []Track) []Track {
	ops, ok := filterFuncs[f.Key]
	if !ok {
		return []Track{}
	}

	ff, ok := ops[f.Op]
	if !ok {
		return []Track{}
	}

	results := []Track{}
	for _, track := range tracks {
		if ff(track, f.Value) {
			results = append(results, track)
		}
	}

	return results
}
