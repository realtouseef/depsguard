package cmd

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/realtouseef/depsguard/internal/baseline"
	"github.com/realtouseef/depsguard/internal/knowledge"
	"github.com/realtouseef/depsguard/internal/parser"
	"github.com/realtouseef/depsguard/internal/selector"
	"github.com/realtouseef/depsguard/internal/util"

	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify dependency explanations before install",
	RunE: func(cmd *cobra.Command, args []string) error {
		base, err := baseline.Load(filepath.Join(util.DepguardDir(), "baseline.json"))
		if err != nil {
			return err
		}
		pkg, err := parser.LoadPackageJSON("package.json")
		if err != nil {
			return err
		}

		current := pkg.AllDependencies()
		newDeps := baseline.NewDependencies(base.Dependencies, current)
		if len(newDeps) == 0 {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "No new dependencies detected.")
			return nil
		}

		allNames := make([]string, 0, len(current))
		for name := range current {
			allNames = append(allNames, name)
		}
		sort.Strings(allNames)

		seed := util.ResolveSeed()
		selected := selector.Select(allNames, 0.4, seed)
		knowledgePath := filepath.Join(util.DepguardDir(), "knowledge.json")
		entries, err := knowledge.Load(knowledgePath)
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

		if len(missing) == 0 {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Dependency knowledge check passed.")
			return nil
		}

		sort.Strings(newDeps)
		sort.Strings(missing)

		builder := strings.Builder{}
		builder.WriteString("DepGuard blocked install.\n")
		builder.WriteString("New dependencies detected:\n")
		for _, dep := range newDeps {
			builder.WriteString("  - ")
			builder.WriteString(dep)
			builder.WriteString("\n")
		}
		builder.WriteString("Missing explanations for selected dependencies:\n")
		for _, dep := range missing {
			builder.WriteString("  - ")
			builder.WriteString(dep)
			builder.WriteString("\n")
		}
		builder.WriteString("Run: depguard explain <dependency-name>\n")
		return fmt.Errorf(builder.String())
	},
}
