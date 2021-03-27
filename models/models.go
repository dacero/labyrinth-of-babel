package models

import (
	"time"
	"strings"
	
	"github.com/russross/blackfriday/v2"
	"github.com/microcosm-cc/bluemonday"
)

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

func (c Cell) HTMLBody() string {
	groomedString := strings.ReplaceAll(c.Body, "\r\n", "\n")
	unsafe := blackfriday.Run([]byte(groomedString))
	output := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	return string(output)
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
	Create_time time.Time
}