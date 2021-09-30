package main

import (
	"github.com/mmcdole/gofeed"
)

func (a *App) compileRSS() error {
	sources, err := a.Models.FeedSource.GetAll()
	if err != nil {
		return err
	}

	fp := gofeed.NewParser()
	for _, source := range sources {
		feed, err := fp.ParseURL(source.URL)
		if err != nil {
			return err
		}
		for _, item := range feed.Items {
			a.Models.FeedItem.Insert(&FeedItem{
				Description: item.Title,
				RelatedLink: &item.Link,
				SourceID: source.ID,
				MediaType: "article",
				PostDate: *item.PublishedParsed,
			})
		}
	}

	return nil
}
