// Package coverage loads a golang coverage report and enrich the entities with the information
package coverage

import (
	"log/slog"
	"path/filepath"
	"strings"

	"golang.org/x/tools/cover"

	"github.com/Fabianexe/go2jenkins/pkg/entity"
)

// LoadCoverage loads the coverage data from the given file
func LoadCoverage(project *entity.Project, coverageReport string, fixMissedLines bool) (*entity.Project, error) {
	profiles, err := cover.ParseProfiles(coverageReport)
	if err != nil {
		return nil, err
	}

	for _, p := range profiles {
		slog.Debug("Profile", "Path", p.FileName, "Blocks", len(p.Blocks))
		found := false
		for _, pack := range project.Packages {
			if !strings.HasPrefix(p.FileName, pack.Name) {
				continue
			}
			found = true
			filename := filepath.Base(p.FileName)
			for _, f := range pack.Files {
				if filepath.Base(f.FilePath) != filename {
					continue
				}
				applyBlocks(f.Methods, p.Blocks, fixMissedLines)
			}

		}
		if !found {
			slog.Warn("Not found source for: " + p.FileName)
		}
	}

	updateLineCoverage(project)
	updateBranchCoverage(project)

	return project, nil
}

func applyBlocks(methods []*entity.Method, blocks []cover.ProfileBlock, fixMissedLines bool) {
	blockLines := make(map[int]struct{}, len(blocks))
	for _, b := range blocks {
		if b.Count == 0 {
			continue
		}
		for _, method := range methods {
			for _, line := range method.Lines {
				if b.StartLine <= line.Number && line.Number <= b.EndLine {
					line.CoverageCount += b.Count
				}
			}
			for _, branch := range method.Branches {
				if b.StartLine <= branch.EndLine && b.EndLine >= branch.StartLine {
					branch.Covered = true
				}
			}
		}
		for i := b.StartLine; i < b.EndLine; i++ {
			blockLines[i] = struct{}{}
		}
	}

	if fixMissedLines {
		for _, method := range methods {
			for _, line := range method.Lines {
				if _, ok := blockLines[line.Number]; !ok {
					line.CoverageCount = determineFromBranches(method.Branches, line.Number)
				}
			}
		}
	}
}

func determineFromBranches(branches []*entity.Branch, number int) int {
	// branches are created by deep first search, so the last branch that cover the line is the one that is interesting
	for i := len(branches) - 1; i >= 0; i-- {
		branch := branches[i]
		if branch.StartLine <= number && number <= branch.EndLine {
			if branch.Covered {
				return 1
			}
			return 0
		}
	}

	return 0
}

func updateLineCoverage(project *entity.Project) {
	for _, pack := range project.Packages {
		for _, f := range pack.Files {
			for _, method := range f.Methods {
				for _, line := range method.Lines {
					isCovered := line.CoverageCount > 0
					method.LineCoverage.AddLine(isCovered)
					f.LineCoverage.AddLine(isCovered)
					pack.LineCoverage.AddLine(isCovered)
					project.LineCoverage.AddLine(isCovered)

				}
			}
		}
	}
}

func updateBranchCoverage(project *entity.Project) {
	for _, pack := range project.Packages {
		for _, f := range pack.Files {
			for _, method := range f.Methods {
				for _, branch := range method.Branches {
					isCovered := branch.Covered
					method.BranchCoverage.AddBranch(isCovered)
					f.BranchCoverage.AddBranch(isCovered)
					pack.BranchCoverage.AddBranch(isCovered)
					project.BranchCoverage.AddBranch(isCovered)

				}
			}
		}
	}
}
