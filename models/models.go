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

func (s Source) String() string {
	return string(s.Source)
}

type CellLink struct {
	Id   string
	Text string
}

type CollectionOfCells struct {
	Name		string
	CellCount	int
}