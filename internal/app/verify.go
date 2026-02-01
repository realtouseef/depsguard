package app

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/realtouseef/depsguard/internal/nodejs"
	"github.com/realtouseef/depsguard/internal/selection"
	"github.com/realtouseef/depsguard/internal/storage"
	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify dependency explanations before install",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := storage.EnsureDepsguardDir(); err != nil {
			return err
		}

		basePath := filepath.Join(storage.DepsguardDir(), "baseline.json")
		base, err := storage.LoadBaseline(basePath)
		if err != nil {
			return err
		}
		pkg, err := nodejs.LoadPackageJSON("package.json")
		if err != nil {
			return err
		}

		current := pkg.GetAllDependencies()
		allNames := make([]string, 0, len(current))
		for name := range current {
			allNames = append(allNames, name)
		}
		sort.Strings(allNames)

		selected := selection.SelectDependencies(allNames, 0.4)
		knowledgePath := filepath.Join(storage.DepsguardDir(), "knowledge.json")
		entries, err := storage.LoadKnowledge(knowledgePath)
		if err != nil {
			return err
		}

		now := time.Now()
		missing := make([]string, 0)
		for _, dep := range selected {
			entry, ok := entries[dep]
			if !ok || !entry.IsValid(now) {
				missing = append(missing, dep)
			}
		}

		newDeps := storage.NewDependencies(base.Dependencies, current)
		sort.Strings(newDeps)
		sort.Strings(missing)

		if len(missing) > 0 {
			builder := strings.Builder{}
			builder.WriteString(color.New(color.FgRed).Sprint("depsguard blocked install."))
			builder.WriteString("\n")
			if len(newDeps) > 0 {
				builder.WriteString("New dependencies detected:\n")
				for _, dep := range newDeps {
					builder.WriteString("  - ")
					builder.WriteString(dep)
					builder.WriteString("\n")
				}
			}
			builder.WriteString("Missing explanations for selected dependencies:\n")
			for _, dep := range missing {
				builder.WriteString("  - ")
				builder.WriteString(dep)
				builder.WriteString("\n")
			}
			builder.WriteString("Run: depsguard explain <dependency-name>\n")
			return fmt.Errorf(builder.String())
		}

		newBase := storage.NewBaseline(current)
		if err := storage.SaveBaseline(basePath, newBase); err != nil {
			return err
		}
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), color.New(color.FgGreen).Sprint("Dependency knowledge check passed."))
		return nil
	},
}
