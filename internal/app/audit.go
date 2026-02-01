package app

import (
	"fmt"
	"path/filepath"
	"sort"
	"time"

	"github.com/fatih/color"
	"github.com/realtouseef/depsguard/internal/nodejs"
	"github.com/realtouseef/depsguard/internal/storage"
	"github.com/spf13/cobra"
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Summarize dependency knowledge coverage",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := storage.EnsureDepsguardDir(); err != nil {
			return err
		}

		pkg, err := nodejs.LoadPackageJSON("package.json")
		if err != nil {
			return err
		}
		deps := pkg.GetAllDependencies()
		knowledgePath := filepath.Join(storage.DepsguardDir(), "knowledge.json")
		entries, err := storage.LoadKnowledge(knowledgePath)
		if err != nil {
			return err
		}

		now := time.Now()
		explained := make([]string, 0)
		unexplained := make([]string, 0)
		soonExpiring := make([]string, 0)
		for dep := range deps {
			entry, ok := entries[dep]
			if ok && entry.IsValid(now) {
				explained = append(explained, dep)
				expiresAt, parseErr := time.Parse(time.RFC3339, entry.ExpiresAt)
				if parseErr == nil && expiresAt.Before(now.Add(30*24*time.Hour)) {
					soonExpiring = append(soonExpiring, dep)
				}
				continue
			}
			unexplained = append(unexplained, dep)
		}

		sort.Strings(explained)
		sort.Strings(unexplained)
		sort.Strings(soonExpiring)

		total := len(deps)
		coverage := 0.0
		if total > 0 {
			coverage = float64(len(explained)) / float64(total) * 100
		}

		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Total dependencies: %d\n", total)
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Explained dependencies: %d (%.1f%%)\n", len(explained), coverage)
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Unexplained dependencies: %d\n", len(unexplained))

		if len(explained) > 0 {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Explained dependencies:")
			for _, dep := range explained {
				entry := entries[dep]
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "  - %s (expires %s)\n", dep, entry.ExpiresAt)
			}
		}

		if len(unexplained) > 0 {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Unexplained dependencies:")
			for _, dep := range unexplained {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", dep)
			}
		}

		if len(soonExpiring) > 0 {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), color.New(color.FgYellow).Sprint("Soon-to-expire explanations (< 30 days):"))
			for _, dep := range soonExpiring {
				entry := entries[dep]
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "  - %s (expires %s)\n", dep, entry.ExpiresAt)
			}
		}

		return nil
	},
}
