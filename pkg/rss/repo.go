package rss

import (
	"encoding/base64"
	"fmt"
	"github.com/SlyMarbo/rss"
	"github.com/linuxsuren/goplay/pkg/data"
	"github.com/linuxsuren/goplay/pkg/util"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func FindMapByKeyword(keyword string) (result map[string]*rss.Feed, titles []string) {
	feeds := FindByKeyword(keyword)
	result = make(map[string]*rss.Feed)
	for i, _ := range feeds {
		feed := feeds[i]
		result[feed.Title] = feed
		titles = append(titles, feed.Title)
	}
	return
}

func loadCache(rssURL string) *rss.Feed {
	cacheDir := os.ExpandEnv("$HOME/.config/goplay/cache")
	_ = os.MkdirAll(cacheDir, 0751)

	rssCacheFile := path.Join(cacheDir, util.HashCodeAsString(rssURL))
	if cacheData, err := ioutil.ReadFile(rssCacheFile); err == nil {
		if feed, err := rss.Parse(cacheData); err != nil {
			_ = os.RemoveAll(rssCacheFile)
		} else {
			return feed
		}
	}
	return nil
}

func saveCache(rssURL string, feed *rss.Feed) {
	cacheDir := os.ExpandEnv("$HOME/.config/goplay/cache")
	_ = os.MkdirAll(cacheDir, 0751)

	rssCacheFile := path.Join(cacheDir, util.HashCodeAsString(rssURL))
	_ = ioutil.WriteFile(rssCacheFile, feed.RawData, 0644)
}

func getRSSByURL(rssURL string) (feed *rss.Feed, err error) {
	if feed = loadCache(rssURL); feed == nil {
		if feed, err = rss.Fetch(rssURL); err == nil {
			saveCache(rssURL, feed)
		}
	}
	return
}

func FindByKeyword(keyword string) (feeds []*rss.Feed) {
	items, _ := data.GetRSSSources()
	feeds = make([]*rss.Feed, 0)
	for _, item := range items {
		rssURL := item.RSS

		if feed, err := getRSSByURL(rssURL); err != nil {
			fmt.Printf("failed to fetch rss feed: %s, error: %v", rssURL, err)
			continue
		} else {
			if containsIgnoreCase(feed.Description, keyword) || containsIgnoreCase(feed.Title, keyword) {
				feeds = append(feeds, feed)
			}
		}
	}
	return
}

func containsIgnoreCase(text, sub string) bool {
	text = strings.ToLower(text)
	sub = strings.ToLower(sub)
	return strings.Contains(text, sub)
}

func ConvertEpisodeMap(items []*rss.Item) (episodes map[string]Episode, titles []string) {
	episodes = make(map[string]Episode)
	for i, _ := range items {
		item := items[i]
		titles = append(titles, item.Title)

		if len(item.Enclosures) > 0 {
			enclosure := item.Enclosures[0]
			episodes[item.Title] = Episode{
				Title:     item.Title,
				AudioLink: enclosure.URL,
				Type:      enclosure.Type,
				Length:    enclosure.Length,
				UID:       base64.StdEncoding.EncodeToString([]byte(item.ID)),
			}
		}
	}
	return
}
