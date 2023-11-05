// Package commands contains all cobra commands that are used from the main
package commands

import (
	"github.com/spf13/cobra"

	"github.com/Fabianexe/go2jenkins/pkg/cleaner"
	"github.com/Fabianexe/go2jenkins/pkg/complexity"
	"github.com/Fabianexe/go2jenkins/pkg/coverage"
	"github.com/Fabianexe/go2jenkins/pkg/source"
	"github.com/Fabianexe/go2jenkins/pkg/writer"
)

func RootCommand() {
	var rootCmd = &cobra.Command{ //nolint:gochecknoglobals
		Use:   "go2jenkins",
		Short: "go2jenkins creates coverage files for golang that the Jenkins coverage plugin can read and display.",
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

			generatedFiles, err := cmd.Flags().GetBool("generatedFiles")
			if err != nil {
				panic(err)
			}

			noneCodeLines, err := cmd.Flags().GetBool("noneCodeLines")
			if err != nil {
				panic(err)
			}

			errorIf, err := cmd.Flags().GetBool("errorIf")
			if err != nil {
				panic(err)
			}

			project, err := source.LoadSources(sourcePath)
			if err != nil {
				panic(err)
			}

			project = cleaner.CleanData(
				project,
				!generatedFiles,
				!noneCodeLines,
				!errorIf,
			)

			project = complexity.AddComplexity(project)

			if coveragePath != "-" {
				project, err = coverage.LoadCoverage(project, coveragePath)

				if err != nil {
					panic(err)
				}
			}

			err = writer.WriteXML(sourcePath, project, outputPath)
			if err != nil {
				panic(err)
			}
		},
	}

	rootCmd.PersistentFlags().StringP(
		"source",
		"s",
		"./",
		"Give the source path to the go project.",
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

	rootCmd.PersistentFlags().Bool(
		"generatedFiles",
		false,
		"If flag is given generated files are part of the output",
	)

	rootCmd.PersistentFlags().Bool(
		"noneCodeLines",
		false,
		"If flag is given non code lines in functions are part of the output",
	)

	rootCmd.PersistentFlags().Bool(
		"errorIf",
		false,
		"If flag is given err ifs are part of the output",
	)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
