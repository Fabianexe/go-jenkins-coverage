// Package coverage loads a golang coverage report and enrich the entities with the information
package coverage

import (
	"path/filepath"
	"strings"

	"golang.org/x/tools/cover"

	"github.com/Fabianexe/go-jenkins-coverage/pkg/entity"
)

// LoadCoverage loads the coverage data from the given file
func LoadCoverage(project *entity.Project, coverageReport string) (*entity.Project, error) {
	profiles, err := cover.ParseProfiles(coverageReport)
	if err != nil {
		return nil, err
	}

	for _, p := range profiles {
		found := false
		for _, pack := range project.Packages {
			if !strings.HasPrefix(p.FileName, pack.Name) {
				continue
			}
			found = true
			filename := filepath.Base(p.FileName)
			for _, f := range pack.Files {
				if f.Name != filename {
					continue
				}
				for _, b := range p.Blocks {
					for _, method := range f.Methods {
						for _, line := range method.Lines {
							if b.StartLine <= line.Number && line.Number <= b.EndLine {
								line.CoverageCount += b.Count
							}
						}
					}
				}
			}

		}
		if !found {
			println("Not found source for :" + p.FileName)
		}
	}

	updateLineCoverage(project)

	return project, nil
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
