// Package commands contains all cobra commands that are used from the main
package commands

import (
	"github.com/spf13/cobra"

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
			packages, err := source.LoadSources(sourcePath)
			if err != nil {
				panic(err)
			}

			err = writer.WriteXML(sourcePath, packages, "coverage.xml")
			if err != nil {
				panic(err)
			}
		},
	}

	rootCmd.PersistentFlags().StringP(
		"source",
		"s",
		".",
		"Apply on this branch if not given is determined by selection.\nIf '`.`' given determine branch from current working dir (if branch is needed but not given branch is tried to determined). ",
	)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
