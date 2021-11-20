package data

import (
	_ "embed"
	"github.com/ghodss/yaml"
)

//go:embed subscription.yaml
var rssResources string

// GetRSSSources returns the pre-defined rss resources
func GetRSSSources() (result []RSSSource, err error) {
	result = make([]RSSSource, 0)
	err = yaml.Unmarshal([]byte(rssResources), &result)
	return
}
