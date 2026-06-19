package command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/inodaf/neoman/cmd/nman/skills"
)

func (c *Command) Skills() {
	err := c.installSkill()
	if err != nil {
		fmt.Printf("neoman: %s.\n", err.Error())
		os.Exit(1)
		return
	}

	os.Exit(0)
}

func (c *Command) installSkill() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ErrSkillHomeDir
	}

	// @TODO: Support other AI Agents.
	skillDir := filepath.Join(homeDir, ".claude", "skills", "nman")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		return ErrSkillCreateDir
	}

	skillPath := filepath.Join(skillDir, "SKILL.md")
	if err := os.WriteFile(skillPath, []byte(skills.NmanSkill), 0644); err != nil {
		return ErrSkillWriteFile
	}

	fmt.Printf("Installed nman skill to %s\n\n", skillPath)
	fmt.Println("Claude can now use nman for documentation lookup.")

	return nil
}

var (
	ErrSkillHomeDir   = fmt.Errorf("Unable to get home directory")
	ErrSkillCreateDir = fmt.Errorf("Unable to create skill directory")
	ErrSkillWriteFile = fmt.Errorf("Unable to write skill file")
)
