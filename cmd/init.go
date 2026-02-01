package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/realtouseef/depsguard/internal/baseline"
	"github.com/realtouseef/depsguard/internal/parser"
	"github.com/realtouseef/depsguard/internal/util"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize depsguard baseline and knowledge store",
	RunE: func(cmd *cobra.Command, args []string) error {
		pkg, err := parser.LoadPackageJSON("package.json")
		if err != nil {
			return err
		}

		deps := pkg.AllDependencies()
		dir := util.DepsguardDir()
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}

		base := baseline.NewBaseline(deps)
		if err := baseline.Save(filepath.Join(dir, "baseline.json"), base); err != nil {
			return err
		}

		knowledgePath := filepath.Join(dir, "knowledge.json")
		if _, err := os.Stat(knowledgePath); os.IsNotExist(err) {
			if err := os.WriteFile(knowledgePath, []byte("{}"), 0o644); err != nil {
				return err
			}
		}

		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "depsguard initialized.")
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Next step: add this to package.json scripts:")
		pretty, _ := json.MarshalIndent(map[string]string{"preinstall": "depsguard verify"}, "", "  ")
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), string(pretty))
		return nil
	},
}
