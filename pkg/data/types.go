package data

type RSSSource struct {
	RSS      string   `yaml:"rss"`
	Category []string `yaml:"category"`
	Language string   `yaml:"language"`
}
