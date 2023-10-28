package entity

import (
	"fmt"
)

type LineCounter struct {
	totalLines   int
	coveredLines int
}

func (c *LineCounter) AddLine(covered bool) {
	c.totalLines++
	if covered {
		c.coveredLines++
	}
}

func (c *LineCounter) String() string {
	return fmt.Sprintf("%.2f", float64(c.coveredLines)/float64(c.totalLines))
}
