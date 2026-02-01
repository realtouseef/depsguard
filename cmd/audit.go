package cmd

import (
	"fmt"
	"path/filepath"
	"sort"
	"time"

	"github.com/realtouseef/depsguard/internal/knowledge"
	"github.com/realtouseef/depsguard/internal/parser"
	"github.com/realtouseef/depsguard/internal/util"

	"github.com/spf13/cobra"
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Summarize dependency knowledge coverage",
	RunE: func(cmd *cobra.Command, args []string) error {
		pkg, err := parser.LoadPackageJSON("package.json")
		if err != nil {
			return err
		}
		deps := pkg.AllDependencies()
		knowledgePath := filepath.Join(util.DepsguardDir(), "knowledge.json")
		entries, err := knowledge.Load(knowledgePath)
		if err != nil {
			return err
		}

		now := time.Now()
		explained := 0
		unexplained := make([]string, 0)
		for dep := range deps {
			if entry, ok := entries[dep]; ok && entry.IsValid(now) {
				explained++
			} else {
				unexplained = append(unexplained, dep)
			}
		}

		sort.Strings(unexplained)
		total := len(deps)
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Total dependencies: %d\n", total)
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Explained dependencies: %d\n", explained)
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Unexplained dependencies: %d\n", total-explained)
		limit := 10
		if len(unexplained) < limit {
			limit = len(unexplained)
		}
		if limit > 0 {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Top unexplained dependencies:")
			for _, dep := range unexplained[:limit] {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", dep)
			}
		}
		return nil
	},
}
