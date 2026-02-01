package git

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

func CurrentUser() string {
	name := readGitConfig("user.name")
	email := readGitConfig("user.email")

	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)

	if name != "" && email != "" {
		return name + " <" + email + ">"
	}
	if name != "" {
		return name
	}
	if email != "" {
		return email
	}
	if value := os.Getenv("USER"); value != "" {
		return value
	}
	if value := os.Getenv("USERNAME"); value != "" {
		return value
	}
	return "unknown"
}

func readGitConfig(key string) string {
	cmd := exec.Command("git", "config", "--get", key)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(bytes.TrimSpace(output))
}
