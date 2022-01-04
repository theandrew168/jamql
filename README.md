# JamQL
Smart playlists for Spotify

## Setup
This project depends on the [Go programming language](https://golang.org/dl/).

## Running
If actively working on frontend templates, set `ENV=dev` to tell the server to reload templates from the filesystem on every page load.
```bash
# make run
ENV=dev go run main.go -conf internal/test/jamql.conf
```

## Testing
```bash
# make test
go test -v ./...
```
