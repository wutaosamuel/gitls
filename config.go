package gitls

// Config is the configuration of program
type Config struct {
	Dir     string   // Dir     -> The directory contain gits
	GitDirs []string // GitDirs -> The git folders' path
}

// NewConfig create a new config
func NewConfig() *Config {
	return &Config{
		Dir:     "",
		GitDirs: make([]string, 0)}
}