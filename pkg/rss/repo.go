package rss

import (
	"encoding/base64"
	"fmt"
	"github.com/SlyMarbo/rss"
	"github.com/linuxsuren/goplay/pkg/data"
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

func FindByKeyword(keyword string) (feeds []*rss.Feed) {
	items, _ := data.GetRSSSources()
	feeds = make([]*rss.Feed, 0)
	for _, item := range items {
		rssURL := item.RSS

		if feed, err := rss.Fetch(rssURL); err != nil {
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
