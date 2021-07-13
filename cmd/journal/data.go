package main

import "time"

type Daily struct {
	Date   time.Time
	Rating float64
}

type Activity struct {
	Date  time.Time
	Major bool
	Desc  string
	Type  string
}
