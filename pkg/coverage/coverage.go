// Package coverage loads a golang coverage report and enrich the entities with the information
package coverage

import (
	"github.com/Fabianexe/go-jenkins-coverage/pkg/entity"
)

// LoadCoverage loads the coverage data from the given file
// TODO: implement
func LoadCoverage(input []*entity.Package, coverageReport string) ([]*entity.Package, error) {
	return input, nil
}
