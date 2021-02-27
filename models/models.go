package models

import "time"

type Cell struct {
	Id          string
	Title       string
	Body        string
	Room        string
	Create_time time.Time
	Update_time time.Time
	Sources     []Source
	Links       []CellLink
}

type Source struct {
	Source string
}

type CellLink struct {
	Id   string
	Text string
}
