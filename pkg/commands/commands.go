// Package commands contains all cobra commands that are used from the main
package commands

import (
	"fmt"
	"log/slog"
	"os"

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
			initLogger()
			slog.Info("Start flag parsing")
			sourcePath, err := cmd.Flags().GetString("source")
			if err != nil {
				slog.Error(fmt.Sprintf("%+v", err))
				os.Exit(1)
			}

			coveragePath, err := cmd.Flags().GetString("coverage")
			if err != nil {
				slog.Error(fmt.Sprintf("%+v", err))
				os.Exit(1)
			}

			outputPath, err := cmd.Flags().GetString("output")
			if err != nil {
				slog.Error(fmt.Sprintf("%+v", err))
				os.Exit(1)
			}

			generatedFiles, err := cmd.Flags().GetBool("generatedFiles")
			if err != nil {
				slog.Error(fmt.Sprintf("%+v", err))
				os.Exit(1)
			}

			noneCodeLines, err := cmd.Flags().GetBool("noneCodeLines")
			if err != nil {
				slog.Error(fmt.Sprintf("%+v", err))
				os.Exit(1)
			}

			errorIf, err := cmd.Flags().GetBool("errorIf")
			if err != nil {
				slog.Error(fmt.Sprintf("%+v", err))
				os.Exit(1)
			}

			cyclomatic, err := cmd.Flags().GetBool("cyclomatic")
			if err != nil {
				slog.Error(fmt.Sprintf("%+v", err))
				os.Exit(1)
			}

			slog.Info("Load sources")
			project, err := source.LoadSources(sourcePath)
			if err != nil {
				slog.Error(fmt.Sprintf("%+v", err))
				os.Exit(1)
			}

			slog.Info("Clean data")
			project = cleaner.CleanData(
				project,
				!generatedFiles,
				!noneCodeLines,
				!errorIf,
			)

			slog.Info("Add complexity")
			project = complexity.AddComplexity(project, cyclomatic, !errorIf)

			if coveragePath != "-" {
				slog.Info("Load coverage")
				project, err = coverage.LoadCoverage(project, coveragePath)

				if err != nil {
					slog.Error(fmt.Sprintf("%+v", err))
					os.Exit(1)
				}
			}

			slog.Info("Write output")
			err = writer.WriteXML(sourcePath, project, outputPath)
			if err != nil {
				slog.Error(fmt.Sprintf("%+v", err))
				os.Exit(1)
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

	rootCmd.PersistentFlags().Bool(
		"cyclomatic",
		false,
		"If flag is given cyclomatic complexity is used instead of cognitive complexity",
	)

	verboseFlag := rootCmd.PersistentFlags().VarPF(
		&verbose,
		"verbose",
		"v",
		"Add verbose output. Multiple -v options increase the verbosity.",
	)
	verboseFlag.NoOptDefVal = "1"

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
