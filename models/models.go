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
	Links       []Cell
}

func (c Cell) HTMLBody() string {
	groomedString := strings.ReplaceAll(c.Body, "\r\n", "\n")
	unsafe := blackfriday.Run([]byte(groomedString))
	output := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	return string(output)
}

func (c Cell) Summary() string {
	if c.Title != "" {
		return c.Title
	} else {
		r := []rune(c.Body)
		if len(r) < 60 {
			return c.Body
		} else {
			return string(r[0:50]) + "..."
		}
	}

}

type Source struct {
	Source string
}

func (s Source) String() string {
	return string(s.Source)
}

type CollectionOfCells struct {
	Name		string
	CellCount	int
	Create_time time.Time
}