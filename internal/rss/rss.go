package rss

import (
	"github.com/dark-warlord14/CVENotifier/internal/errors"
	"github.com/mmcdole/gofeed"
)

func ParseFeed(feedURL string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		return nil, &errors.RSSFeedError{Message: "Failed to parse RSS feed: " + err.Error()}
	}

	return feed, nil
}
