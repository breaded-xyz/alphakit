package main

// Values for build info set at build time with ldflags
var (
	buildGitTag    string
	buildGitCommit string
	buildTime      string
	buildUser      string
)
