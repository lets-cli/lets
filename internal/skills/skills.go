package skills

import _ "embed"

const (
	LetsName     = "lets"
	SkillFile    = "SKILL.md"
	SkillsRelDir = ".agents/skills"
)

//go:embed lets/SKILL.md
var letsSkill []byte

func LetsSkill() []byte {
	content := make([]byte, len(letsSkill))
	copy(content, letsSkill)

	return content
}
