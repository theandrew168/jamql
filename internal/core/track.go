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

var sampleTracks = []Track{
	Track{
		ID:      "1",
		Name:    "Jesusland",
		Artist:  "Ben Folds",
		Album:   "Songs For Silverman",
		Artwork: "https://bulma.io/images/placeholders/64x64.png",
		Genre:   "Alt Rock",
		Year:    "2005",
	},
	Track{
		ID:      "2",
		Name:    "Landed",
		Artist:  "Ben Folds",
		Album:   "Songs For Silverman",
		Artwork: "https://bulma.io/images/placeholders/64x64.png",
		Genre:   "Alt Rock",
		Year:    "2005",
	},
	Track{
		ID:      "3",
		Name:    "Late",
		Artist:  "Ben Folds",
		Album:   "Songs For Silverman",
		Artwork: "https://bulma.io/images/placeholders/64x64.png",
		Genre:   "Alt Rock",
		Year:    "2005",
	},
	Track{
		ID:      "4",
		Name:    "Prison Food",
		Artist:  "Ben Folds",
		Album:   "Songs For Silverman",
		Artwork: "https://bulma.io/images/placeholders/64x64.png",
		Genre:   "Alt Rock",
		Year:    "2005",
	},
	Track{
		ID:      "5",
		Name:    "One Angry Dwarf and 200 Solemn Faces",
		Artist:  "Ben Folds Five",
		Album:   "Whatever and Ever Amen",
		Artwork: "https://bulma.io/images/placeholders/64x64.png",
		Genre:   "Alt Rock",
		Year:    "1997",
	},
	Track{
		ID:      "6",
		Name:    "Brick",
		Artist:  "Ben Folds Five",
		Album:   "Whatever and Ever Amen",
		Artwork: "https://bulma.io/images/placeholders/64x64.png",
		Genre:   "Alt Rock",
		Year:    "1997",
	},
	Track{
		ID:      "7",
		Name:    "Battle of Who Could Care Less",
		Artist:  "Ben Folds Five",
		Album:   "Whatever and Ever Amen",
		Artwork: "https://bulma.io/images/placeholders/64x64.png",
		Genre:   "Alt Rock",
		Year:    "1997",
	},
	Track{
		ID:      "8",
		Name:    "Teflon",
		Artist:  "The Mars Volta",
		Album:   "Octahedron",
		Artwork: "https://bulma.io/images/placeholders/64x64.png",
		Genre:   "Prog Rock",
		Year:    "2009",
	},
	Track{
		ID:      "9",
		Name:    "Cotopaxi",
		Artist:  "The Mars Volta",
		Album:   "Octahedron",
		Artwork: "https://bulma.io/images/placeholders/64x64.png",
		Genre:   "Prog Rock",
		Year:    "2009",
	},
	Track{
		ID:      "10",
		Name:    "Desperate Graves",
		Artist:  "The Mars Volta",
		Album:   "Octahedron",
		Artwork: "https://bulma.io/images/placeholders/64x64.png",
		Genre:   "Prog Rock",
		Year:    "2009",
	},
}
