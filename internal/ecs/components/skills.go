package components

// SkillType represents different types of skills
type SkillType int

const (
	SkillTypePassive SkillType = iota // Permanent stat bonuses
	SkillTypeActive                   // Usable abilities in combat
	SkillTypeTrait                    // Special character traits
	SkillTypeUpgrade                  // Equipment or item upgrades
)

func (st SkillType) String() string {
	switch st {
	case SkillTypePassive:
		return "Passive"
	case SkillTypeActive:
		return "Active"
	case SkillTypeTrait:
		return "Trait"
	case SkillTypeUpgrade:
		return "Upgrade"
	default:
		return "Unknown"
	}
}

// SkillEffect represents the effect a skill has when learned
type SkillEffect struct {
	Type        string      // "stat_bonus", "ability_unlock", "passive_effect"
	Target      string      // What is affected (HP, Attack, etc.)
	Value       int         // Numeric value for the effect
	Description string      // Human readable description
	Data        interface{} // Additional effect data
}

// Skill represents a learnable skill in the game
type Skill struct {
	ID            string        // Unique skill identifier
	Name          string        // Display name
	Description   string        // Detailed description
	Type          SkillType     // Type of skill
	JobClass      JobType       // Which job can learn this skill
	Tier          int           // Skill tier/level (1-5)
	Prerequisites []string      // Required skill IDs to unlock this skill
	SkillPoints   int           // Cost in skill points to learn
	Effects       []SkillEffect // What this skill does when learned
	IconPath      string        // Path to skill icon
	IsLearned     bool          // Whether the character has learned this skill
	IsAvailable   bool          // Whether the skill can be learned now
}

// SkillNode represents a skill in the visual skill tree
type SkillNode struct {
	Skill    *Skill  // The actual skill data
	X, Y     int     // Position in the skill tree grid
	Children []string // Child skill IDs that this unlocks
}

// SkillTree represents a complete skill tree for a job class
type SkillTree struct {
	JobClass JobType                // Which job this tree belongs to
	Name     string                 // Display name of the tree
	Nodes    map[string]*SkillNode  // All skills in this tree, keyed by skill ID
	Layout   [][]*SkillNode         // 2D grid layout for visual representation
	MaxTier  int                    // Maximum tier available in this tree
}

// SkillsComponent manages a character's learned skills and available skill points
type SkillsComponent struct {
	AvailablePoints int                    // Unspent skill points
	TotalPoints     int                    // Total skill points ever earned
	LearnedSkills   map[string]*Skill      // Skills the character has learned
	SkillTrees      map[JobType]*SkillTree // Available skill trees for this character
	ActiveAbilities []string               // Currently equipped active skills
	MaxActiveSlots  int                    // Maximum active abilities that can be equipped
}

// NewSkillsComponent creates a new skills component for a character
func NewSkillsComponent(jobClass JobType) *SkillsComponent {
	return &SkillsComponent{
		AvailablePoints: 0,
		TotalPoints:     0,
		LearnedSkills:   make(map[string]*Skill),
		SkillTrees:      make(map[JobType]*SkillTree),
		ActiveAbilities: make([]string, 0),
		MaxActiveSlots:  4, // Default 4 active ability slots
	}
}

// CanLearnSkill checks if a skill can be learned by this character
func (sc *SkillsComponent) CanLearnSkill(skill *Skill) bool {
	// Check if already learned
	if _, learned := sc.LearnedSkills[skill.ID]; learned {
		return false
	}

	// Check if character has enough skill points
	if sc.AvailablePoints < skill.SkillPoints {
		return false
	}

	// Check prerequisites
	for _, prereqID := range skill.Prerequisites {
		if _, learned := sc.LearnedSkills[prereqID]; !learned {
			return false
		}
	}

	return true
}

// LearnSkill learns a skill and applies its effects
func (sc *SkillsComponent) LearnSkill(skill *Skill) bool {
	if !sc.CanLearnSkill(skill) {
		return false
	}

	// Spend skill points
	sc.AvailablePoints -= skill.SkillPoints

	// Learn the skill
	learnedSkill := *skill // Copy the skill
	learnedSkill.IsLearned = true
	sc.LearnedSkills[skill.ID] = &learnedSkill

	return true
}

// AddSkillPoints adds skill points to the character
func (sc *SkillsComponent) AddSkillPoints(points int) {
	sc.AvailablePoints += points
	sc.TotalPoints += points
}

// GetLearnedSkillsByType returns all learned skills of a specific type
func (sc *SkillsComponent) GetLearnedSkillsByType(skillType SkillType) []*Skill {
	var skills []*Skill
	for _, skill := range sc.LearnedSkills {
		if skill.Type == skillType {
			skills = append(skills, skill)
		}
	}
	return skills
}

// GetActiveAbilities returns currently equipped active abilities
func (sc *SkillsComponent) GetActiveAbilities() []*Skill {
	var abilities []*Skill
	for _, skillID := range sc.ActiveAbilities {
		if skill, exists := sc.LearnedSkills[skillID]; exists && skill.Type == SkillTypeActive {
			abilities = append(abilities, skill)
		}
	}
	return abilities
}

// EquipActiveAbility equips an active ability if there's a free slot
func (sc *SkillsComponent) EquipActiveAbility(skillID string) bool {
	// Check if skill is learned and is active type
	skill, learned := sc.LearnedSkills[skillID]
	if !learned || skill.Type != SkillTypeActive {
		return false
	}

	// Check if already equipped
	for _, equipped := range sc.ActiveAbilities {
		if equipped == skillID {
			return false
		}
	}

	// Check if we have a free slot
	if len(sc.ActiveAbilities) >= sc.MaxActiveSlots {
		return false
	}

	sc.ActiveAbilities = append(sc.ActiveAbilities, skillID)
	return true
}

// UnequipActiveAbility removes an active ability from equipped slots
func (sc *SkillsComponent) UnequipActiveAbility(skillID string) bool {
	for i, equipped := range sc.ActiveAbilities {
		if equipped == skillID {
			// Remove from slice
			sc.ActiveAbilities = append(sc.ActiveAbilities[:i], sc.ActiveAbilities[i+1:]...)
			return true
		}
	}
	return false
}