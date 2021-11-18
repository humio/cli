package main

// GoReleaser will override these when building: https://goreleaser.com/customization/build/
var (
	commit  = "none"
	date    = "unknown"
	version = "master"
)

func SetVersion(v, c, d string) {
	version = v
	commit = c
	date = d
}
