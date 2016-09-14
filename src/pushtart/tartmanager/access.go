package tartmanager

import (
	"pushtart/config"
)

func Exists(pushURL string) bool {
	if config.All().Tarts == nil {
		config.All().Tarts = map[string]config.Tart{}
		return false
	}
	if _, ok := config.All().Tarts[pushURL]; ok {
		return true
	}
	return false
}

func Get(pushURL string) config.Tart {
	if config.All().Tarts == nil {
		return config.Tart{}
	}
	return config.All().Tarts[pushURL]
}

func New(pushURL, owner string) {
	if config.All().Tarts == nil {
		config.All().Tarts = map[string]config.Tart{}
	}
	config.All().Tarts[pushURL] = config.Tart{
		Name:      pushURL,
		PushURL:   pushURL,
		IsRunning: false,
		Owner:     owner,
		PID:       -1,
	}
	config.Flush()
}
