package app

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/realtouseef/depsguard/internal/git"
	"github.com/realtouseef/depsguard/internal/nodejs"
	"github.com/realtouseef/depsguard/internal/storage"
	"github.com/spf13/cobra"
)

var explainCmd = &cobra.Command{
	Use:   "explain <dependency-name>",
	Short: "Explain a dependency and store knowledge",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dep := strings.TrimSpace(args[0])
		if dep == "" {
			return fmt.Errorf("dependency name cannot be empty")
		}

		pkg, err := nodejs.LoadPackageJSON("package.json")
		if err != nil {
			return err
		}
		if _, ok := pkg.GetAllDependencies()[dep]; !ok {
			return fmt.Errorf("dependency %q not found in package.json", dep)
		}

		if err := storage.EnsureDepsguardDir(); err != nil {
			return err
		}

		var summary string
		if err := survey.AskOne(&survey.Input{
			Message: "Summary (what does it do?)",
		}, &summary); err != nil {
			return err
		}
		summary = strings.TrimSpace(summary)
		if summary == "" {
			return fmt.Errorf("summary cannot be empty")
		}

		defaultUser := git.CurrentUser()
		explainedBy := defaultUser
		if err := survey.AskOne(&survey.Input{
			Message: "Explained by",
			Default: defaultUser,
		}, &explainedBy); err != nil {
			return err
		}
		explainedBy = strings.TrimSpace(explainedBy)
		if explainedBy == "" {
			return fmt.Errorf("explained by cannot be empty")
		}

		entry := storage.Entry{
			Summary:     summary,
			ExplainedBy: explainedBy,
			ExpiresAt:   time.Now().Add(90 * 24 * time.Hour).Format(time.RFC3339),
		}
		knowledgePath := filepath.Join(storage.DepsguardDir(), "knowledge.json")
		entries, err := storage.LoadKnowledge(knowledgePath)
		if err != nil {
			return err
		}
		entries[dep] = entry
		if err := storage.SaveKnowledge(knowledgePath, entries); err != nil {
			return err
		}

		_, _ = fmt.Fprintln(cmd.OutOrStdout(), color.New(color.FgGreen).Sprint("Explanation saved."))
		return nil
	},
}
