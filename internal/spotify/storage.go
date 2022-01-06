package spotify

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/golangcollege/sessions"
	"github.com/zmb3/spotify/v2"

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
	// early exit if no auth token is present
	token := s.session.GetString(r, "token")
	if token == "" {
		return nil, core.ErrUnauthorized
	}

	// handle no filters as a special case (return no tracks)
	if len(filters) == 0 {
		return nil, nil
	}

	// just in case the filters don't build to anything (return no tracks)
	query := NewQuery(filters)
	if query.String() == "" {
		return nil, nil
	}

	client := spotify.New(NewAccessTokenClient(token))
	result, err := client.Search(
		context.Background(),
		query.String(),
		spotify.SearchTypeTrack,
		spotify.Limit(50),
	)
	if err != nil {
		if err.Error() == "Invalid access token" {
			return nil, core.ErrUnauthorized
		}
		return nil, err
	}

	// translate API results into []core.Track
	var tracks []core.Track
	for _, t := range result.Tracks.Tracks {
		year := 0

		// check for full date first
		release, err := time.Parse("2006-01-02", t.Album.ReleaseDate)
		if err == nil {
			year = release.Year()
		// otherwise just check for a year
		} else {
			releaseYear, err := strconv.Atoi(t.Album.ReleaseDate)
			if err == nil {
				year = releaseYear
			}
		}

		// default to placeholder artwork
		artwork := "https://bulma.io/images/placeholders/64x64.png"
		for _, image := range t.Album.Images {
			if image.Width == 64 && image.Height == 64 {
				artwork = image.URL
				break
			}
		}

		track := core.Track{
			ID:      t.ID.String(),
			Name:    t.Name,
			Artist:  t.Artists[0].Name,
			Album:   t.Album.Name,
			Artwork: artwork,
			Genre:   query.Genre,
			Year:    strconv.Itoa(year),
		}
		tracks = append(tracks, track)
	}

	// apply each filter
	for _, filter := range filters {
		tracks = filter.Apply(tracks)
	}

	return tracks, nil
}

func (s *storage) SaveTracks(r *http.Request, tracks []core.Track, name, desc string) error {
	// early exit if no auth token is present
	token := s.session.GetString(r, "token")
	if token == "" {
		return core.ErrUnauthorized
	}

	// grab the current user
	client := spotify.New(NewAccessTokenClient(token))
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		if err.Error() == "Invalid access token" {
			return core.ErrUnauthorized
		}
		return err
	}

	// create the playlist
	playlist, err := client.CreatePlaylistForUser(
		context.Background(),
		user.ID,
		name,
		desc,
		true,
		false,
	)
	if err != nil {
		if err.Error() == "Invalid access token" {
			return core.ErrUnauthorized
		}
		return err
	}

	// build a slice of track IDs
	var ids []spotify.ID
	for _, track := range tracks {
		ids = append(ids, spotify.ID(track.ID))
	}

	// add tracks to the playlist
	_, err = client.AddTracksToPlaylist(
		context.Background(),
		playlist.ID,
		ids...,
	)
	if err != nil {
		if err.Error() == "Invalid access token" {
			return core.ErrUnauthorized
		}
		return err
	}

	return nil
}

// hack to get an access token to the spotify library
// why is this so difficult... the lib coupled itself too
//   strongly to handling the OAuth flow, IMO
type transport struct {
	token string
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer " + t.token)
	return http.DefaultTransport.RoundTrip(req)
}

func NewAccessTokenClient(token string) *http.Client {
	client := http.Client{
		Transport: &transport{token: token},
	}
	return &client
}
