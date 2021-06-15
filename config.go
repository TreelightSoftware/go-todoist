package todoist

import "os"

type Configuration struct {
	AuthToken string // should be set if, and only if, you are using this for a single user
}

var config *Configuration

func setup() {
	if config != nil {
		return
	}
	config = &Configuration{}
	config.AuthToken = envHelper("TODOIST_AUTH_TOKEN", "")
}

func envHelper(key, defaultValue string) string {
	found := os.Getenv(key)
	if found == "" {
		found = defaultValue
	}
	return found
}

func init() {
	setup()
}
