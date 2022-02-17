package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/mmcdole/gofeed/rss"
)

type CustomRSSTranslator struct {
	defaultTranslator *gofeed.DefaultRSSTranslator
}

func NewCustomRSSTranslator() *CustomRSSTranslator {
	t := &CustomRSSTranslator{}
	t.defaultTranslator = &gofeed.DefaultRSSTranslator{}
	return t
}

func (ct *CustomRSSTranslator) Translate(feed interface{}) (*gofeed.Feed, error) {
	rss, found := feed.(*rss.Feed)
	if !found {
		return nil, fmt.Errorf("Feed did not match expected type of *rss.Feed")
	}
	f, err := ct.defaultTranslator.Translate(rss)
	if err != nil {
		return nil, err
	}
	for i, item := range rss.Items {
		f.Items[i].Custom = map[string]string{"Comments": item.Comments}
	}
	return f, nil
}
