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
	if c.totalLines == 0 {
		return "1.00"
	}

	return fmt.Sprintf("%.2f", float64(c.coveredLines)/float64(c.totalLines))
}

type BranchCounter struct {
	totalBranches   int
	coveredBranches int
}

func (c *BranchCounter) AddBranch(covered bool) {
	c.totalBranches++
	if covered {
		c.coveredBranches++
	}
}

func (c *BranchCounter) String() string {
	if c.totalBranches == 0 {
		return "1.00"
	}

	return fmt.Sprintf("%.2f", float64(c.coveredBranches)/float64(c.totalBranches))
}
