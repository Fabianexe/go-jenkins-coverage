// Package commands contains all cobra commands that are used from the main
package commands

import (
	"github.com/spf13/cobra"

	"github.com/Fabianexe/go-jenkins-coverage/pkg/cleaner"
	"github.com/Fabianexe/go-jenkins-coverage/pkg/complexity"
	"github.com/Fabianexe/go-jenkins-coverage/pkg/coverage"
	"github.com/Fabianexe/go-jenkins-coverage/pkg/source"
	"github.com/Fabianexe/go-jenkins-coverage/pkg/writer"
)

func RootCommand() {
	var rootCmd = &cobra.Command{ //nolint:gochecknoglobals
		Use:   "jcoverage",
		Short: "jcoverage creates coverage files for golang that the Jenkins coverage plugin can read and display.",
		Run: func(cmd *cobra.Command, _ []string) {
			sourcePath, err := cmd.Flags().GetString("source")
			if err != nil {
				panic(err)
			}
			coveragePath, err := cmd.Flags().GetString("coverage")
			if err != nil {
				panic(err)
			}
			outputPath, err := cmd.Flags().GetString("output")
			if err != nil {
				panic(err)
			}

			packages, err := source.LoadSources(sourcePath)
			if err != nil {
				panic(err)
			}

			packages = cleaner.CleanData(packages)

			packages = complexity.AddComplexity(packages)

			if coveragePath == "-" {
				packages, err = coverage.LoadCoverage(packages, coveragePath)
			}

			if err != nil {
				panic(err)
			}

			err = writer.WriteXML(sourcePath, packages, outputPath)
			if err != nil {
				panic(err)
			}
		},
	}

	rootCmd.PersistentFlags().StringP(
		"source",
		"s",
		".",
		"Give The source path to the go project.",
	)

	rootCmd.PersistentFlags().StringP(
		"coverage",
		"c",
		"-",
		"Give the path to the coverage.out file. If omitted no coverage data is considered",
	)

	rootCmd.PersistentFlags().StringP(
		"output",
		"o",
		"coverage.xml",
		"The output file name",
	)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
