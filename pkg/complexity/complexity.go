// Package complexity enriches the enties with complexity metrics
package complexity

import (
	"log/slog"

	"github.com/Fabianexe/go2jenkins/pkg/entity"
)

// AddComplexity adds complexity metrics to the packages
func AddComplexity(project *entity.Project, useCyclomaticComplexity bool, ignoreErrorIF bool) *entity.Project {
	if useCyclomaticComplexity {
		slog.Info("Use cyclomatic complexity")
	} else {
		slog.Info("Use cognitive complexity")
	}
	for _, p := range project.Packages {
		for _, f := range p.Files {
			for _, method := range f.Methods {
				if useCyclomaticComplexity {
					method.Complexity = getCyclomaticComplexity(method.Body, ignoreErrorIF)
				} else {
					method.Complexity = getCognitiveComplexity(method.Body, ignoreErrorIF)
				}
			}
		}
	}

	return project
}
