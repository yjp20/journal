package main

import (
	"os/exec"
	"strings"
	"time"

	readability "github.com/go-shiori/go-readability"
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
				SourceID:    source.ID,
				MediaType:   "article",
				PostDate:    *item.PublishedParsed,
			})
		}
	}

	return nil
}

func (a *App) syncReadableMedia() error {
	media, err := a.Models.Media.GetAll()
	if err != nil {
		return err
	}
	for _, item := range media {
		if item.Cart && item.MediaType == "article" {
			link := item.RelatedLink
			article, err := readability.FromURL(*link, 30*time.Second)
			if err != nil {
				continue
			}
			cmd := exec.Command("pandoc", "-o", item.Description + ".pdf")
			cmd.Stdin = strings.NewReader(article.Content)
			err = cmd.Run()
			if err != nil {
				println(err)
				continue
			}
		}
	}
	return nil
}
