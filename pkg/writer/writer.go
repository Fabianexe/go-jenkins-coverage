// Package writer writes xml file based on the cobertura dtd:
// https://github.com/cobertura/cobertura/blob/master/cobertura/src/test/resources/dtds/coverage-04.dtd
package writer

import (
	"encoding/xml"
	"log/slog"
	"os"

	"github.com/Fabianexe/go2jenkins/pkg/entity"
)

func WriteXML(path string, project *entity.Project, outPath string) error {
	xmlCoverage := ConvertToCobertura(path, project)

	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}

	encoder := xml.NewEncoder(outFile)
	encoder.Indent("", "\t")

	slog.Info("Write coverage to file", "Path", outPath)
	err = encoder.Encode(xmlCoverage)
	if err != nil {
		return err
	}

	if err := outFile.Close(); err != nil {
		return err
	}

	return nil
}
