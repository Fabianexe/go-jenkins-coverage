package writer

import (
	"strconv"

	"github.com/Fabianexe/go-jenkins-coverage/pkg/entity"
)

func ConvertToCobertura(path string, project *entity.Project) *Coverage {
	pkgs := project.Packages
	coverage := &Coverage{
		Sources: &Sources{
			Sources: []*Source{
				{
					Path: path,
				},
			},
		},
		LineRate:   project.LineCoverage.String(),
		BranchRate: project.BranchCoverage.String(),
	}

	packages := &Packages{
		Packages: make([]*Package, 0, len(pkgs)),
	}
	for _, pkg := range pkgs {
		packageCov := &Package{
			Name:       pkg.Name,
			LineRate:   pkg.LineCoverage.String(),
			BranchRate: pkg.BranchCoverage.String(),
		}

		classes := &Classes{
			Classes: make([]*Class, 0, len(pkg.Files)),
		}

		for _, file := range pkg.Files {
			class := &Class{
				Name:       file.Name,
				Filename:   file.FilePath,
				LineRate:   file.LineCoverage.String(),
				BranchRate: file.BranchCoverage.String(),
			}

			mmethods := &Methods{
				Methods: make([]*Method, 0, len(file.Methods)),
			}

			classLines := &Lines{
				Lines: make([]*Line, 0, 1024),
			}

			for _, method := range file.Methods {
				xmlMethod := &Method{
					Name:       method.Name,
					LineRate:   method.LineCoverage.String(),
					BranchRate: method.BranchCoverage.String(),
				}

				methodsLines := &Lines{
					Lines: make([]*Line, 0, len(method.Lines)),
				}

				for _, line := range method.Lines {
					xmlLine := &Line{
						Number: strconv.Itoa(line.Number),
						Hits:   strconv.Itoa(line.CoverageCount),
					}
					methodsLines.Lines = append(methodsLines.Lines, xmlLine)
					classLines.Lines = append(classLines.Lines, xmlLine)
				}

				if len(methodsLines.Lines) != 0 {
					xmlMethod.Lines = methodsLines
				}

				mmethods.Methods = append(mmethods.Methods, xmlMethod)
			}
			if len(mmethods.Methods) != 0 {
				class.Methods = mmethods
			}

			if len(classLines.Lines) != 0 {
				class.Lines = classLines
			}

			classes.Classes = append(classes.Classes, class)
		}

		if len(classes.Classes) != 0 {
			packageCov.Classes = classes
		}

		packages.Packages = append(packages.Packages, packageCov)
	}
	if len(packages.Packages) != 0 {
		coverage.Packages = packages
	}

	return coverage
}
