// Package cleaner cleans the coverage data and discard lines, classes and packages that are not relevant
package cleaner

import (
	"github.com/Fabianexe/go2jenkins/pkg/entity"
)

// CleanData cleanes the package data
func CleanData(
	project *entity.Project,
	cGeneratedFiles bool,
	cNoneCodeLines bool,
	cErrorIf bool,
) *entity.Project {
	if cGeneratedFiles {
		project = cleanGeneratedFiles(project)
	}

	if cNoneCodeLines {
		project = cleanNoneCodeLines(project)
	}

	if cErrorIf {
		project = cleanErrorIf(project)
	}

	return project
}
