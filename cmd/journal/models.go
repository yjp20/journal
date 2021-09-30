package main

import (
	"errors"
)

var ErrRecordNotFound = errors.New("record not found")

type Models struct {
	Media      *MediaModel
	Todo       *TodoModel
	FeedSource *FeedSourceModel
	FeedItem   *FeedItemModel
}
