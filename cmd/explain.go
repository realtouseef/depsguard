package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/realtouseef/depsguard/internal/knowledge"
	"github.com/realtouseef/depsguard/internal/util"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type explainModel struct {
	prompts []string
	answers []string
	index   int
	input   textinput.Model
	err     error
	done    bool
}

func newExplainModel() explainModel {
	input := textinput.New()
	input.Focus()
	return explainModel{
		prompts: []string{
			"What does this dependency do?",
			"Where is it used?",
			"Could it be removed?",
		},
		answers: make([]string, 0, 3),
		input:   input,
	}
}

func (m explainModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m explainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.err = errors.New("explanation cancelled")
			return m, tea.Quit
		case tea.KeyEnter:
			m.answers = append(m.answers, strings.TrimSpace(m.input.Value()))
			m.index++
			if m.index >= len(m.prompts) {
				m.done = true
				return m, tea.Quit
			}
			m.input.SetValue("")
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m explainModel) View() string {
	if m.done {
		return "Saving explanation...\n"
	}
	if m.err != nil {
		return "Explanation cancelled.\n"
	}
	return fmt.Sprintf("%s\n%s\n\nPress Enter to continue.", m.prompts[m.index], m.input.View())
}

var explainCmd = &cobra.Command{
	Use:   "explain <dependency-name>",
	Short: "Explain a dependency and store knowledge",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dep := strings.TrimSpace(args[0])
		if dep == "" {
			return fmt.Errorf("dependency name cannot be empty")
		}

		if err := os.MkdirAll(util.DepsguardDir(), 0o755); err != nil {
			return err
		}

		program := tea.NewProgram(newExplainModel())
		result, err := program.Run()
		if err != nil {
			return err
		}
		model := result.(explainModel)
		if model.err != nil {
			return model.err
		}
		if len(model.answers) != 3 {
			return fmt.Errorf("explanation incomplete")
		}

		summary := fmt.Sprintf("What: %s\nWhere: %s\nRemoval: %s", model.answers[0], model.answers[1], model.answers[2])
		entry := knowledge.Entry{
			Summary:     summary,
			ExplainedBy: util.CurrentUser(),
			ExpiresAt:   time.Now().Add(90 * 24 * time.Hour).Format(time.RFC3339),
		}
		knowledgePath := filepath.Join(util.DepsguardDir(), "knowledge.json")
		entries, err := knowledge.Load(knowledgePath)
		if err != nil {
			return err
		}
		entries[dep] = entry
		if err := knowledge.Save(knowledgePath, entries); err != nil {
			return err
		}

		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Explanation saved.")
		return nil
	},
}
