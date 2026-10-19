// Package version holds the shortlink build version.
package version

// Version is the shortlink release version, overridden at build time with
// -ldflags "-X github.com/BotJakeDemo/shortlink/internal/version.Version=...".
var Version = "0.0.0-dev"
