package version

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func String() string {
	return "mend " + Version + " (" + Commit + ", " + Date + ")"
}
