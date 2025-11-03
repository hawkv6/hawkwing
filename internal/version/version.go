package version

// Version is set via ldflags during build
var Version = "dev"

func GetVersion() string {
	return Version
}
