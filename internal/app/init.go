package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/realtouseef/depsguard/internal/nodejs"
	"github.com/realtouseef/depsguard/internal/storage"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize depsguard baseline and knowledge store",
	RunE: func(cmd *cobra.Command, args []string) error {
		pkg, err := nodejs.LoadPackageJSON("package.json")
		if err != nil {
			return err
		}

		deps := pkg.GetAllDependencies()
		dir := storage.DepsguardDir()
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}

		base := storage.NewBaseline(deps)
		if err := storage.SaveBaseline(filepath.Join(dir, "baseline.json"), base); err != nil {
			return err
		}

		if err := storage.SaveKnowledge(filepath.Join(dir, "knowledge.json"), map[string]storage.Entry{}); err != nil {
			return err
		}

		if err := ensureGitignoreHasDepsguard(); err != nil {
			return err
		}

		_, _ = fmt.Fprintln(cmd.OutOrStdout(), color.New(color.FgGreen).Sprint("depsguard initialized."))
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Next step: add this to package.json scripts:")
		pretty, _ := json.MarshalIndent(map[string]string{"preinstall": "depsguard verify"}, "", "  ")
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), string(pretty))
		return nil
	},
}

func ensureGitignoreHasDepsguard() error {
	const entry = ".depsguard/"
	data, err := os.ReadFile(".gitignore")
	if err != nil {
		if os.IsNotExist(err) {
			return os.WriteFile(".gitignore", []byte(entry+"\n"), 0o644)
		}
		return err
	}
	contents := string(data)
	for _, line := range strings.Split(contents, "\n") {
		if strings.TrimSpace(line) == entry {
			return nil
		}
	}
	builder := strings.Builder{}
	builder.WriteString(contents)
	if !strings.HasSuffix(contents, "\n") && contents != "" {
		builder.WriteString("\n")
	}
	builder.WriteString(entry)
	builder.WriteString("\n")
	return os.WriteFile(".gitignore", []byte(builder.String()), 0o644)
}
