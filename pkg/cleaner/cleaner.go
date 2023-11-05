// Package cleaner cleans the coverage data and discard lines, classes and packages that are not relevant
package cleaner

import (
	"github.com/Fabianexe/go-jenkins-coverage/pkg/entity"
)

// CleanData cleanes the package data
func CleanData(project *entity.Project) *entity.Project {
	project = cleanGeneratedFiles(project)
	project = cleanNoneCodeLines(project)
	// todo clean error ifs (optional)

	return project
}
