package skills

import (
	"fmt"

	"github.com/jrecuero/myrpg/internal/ecs/components"
)

// SkillRegistry manages all available skills in the game
type SkillRegistry struct {
	skills     map[string]*components.Skill                 // All skills keyed by ID
	skillTrees map[components.JobType]*components.SkillTree // Skill trees by job class
}

// NewSkillRegistry creates a new skill registry
func NewSkillRegistry() *SkillRegistry {
	registry := &SkillRegistry{
		skills:     make(map[string]*components.Skill),
		skillTrees: make(map[components.JobType]*components.SkillTree),
	}

	// Initialize with default skills
	registry.initializeDefaultSkills()
	return registry
}

// RegisterSkill adds a skill to the registry
func (sr *SkillRegistry) RegisterSkill(skill *components.Skill) error {
	if skill.ID == "" {
		return fmt.Errorf("skill ID cannot be empty")
	}

	if _, exists := sr.skills[skill.ID]; exists {
		return fmt.Errorf("skill with ID %s already exists", skill.ID)
	}

	sr.skills[skill.ID] = skill
	return nil
}

// GetSkill retrieves a skill by ID
func (sr *SkillRegistry) GetSkill(skillID string) (*components.Skill, bool) {
	skill, exists := sr.skills[skillID]
	return skill, exists
}

// GetSkillsByJob returns all skills for a specific job class
func (sr *SkillRegistry) GetSkillsByJob(jobClass components.JobType) []*components.Skill {
	var jobSkills []*components.Skill
	for _, skill := range sr.skills {
		if skill.JobClass == jobClass {
			jobSkills = append(jobSkills, skill)
		}
	}
	return jobSkills
}

// GetSkillTree returns the skill tree for a job class
func (sr *SkillRegistry) GetSkillTree(jobClass components.JobType) (*components.SkillTree, bool) {
	tree, exists := sr.skillTrees[jobClass]
	return tree, exists
}

// initializeDefaultSkills creates the default skill trees for each job class
func (sr *SkillRegistry) initializeDefaultSkills() {
	// Initialize Warrior skills
	sr.initializeWarriorSkills()

	// Initialize Mage skills
	sr.initializeMageSkills()

	// Initialize Rogue skills
	sr.initializeRogueSkills()

	// Initialize Cleric skills
	sr.initializeClericSkills()

	// Initialize Archer skills
	sr.initializeArcherSkills()
}

// initializeWarriorSkills creates the warrior skill tree
func (sr *SkillRegistry) initializeWarriorSkills() {
	// Tier 1 Skills
	toughSkin := &components.Skill{
		ID:            "warrior_tough_skin",
		Name:          "Tough Skin",
		Description:   "Increases maximum HP by 10 points.",
		Type:          components.SkillTypePassive,
		JobClass:      components.JobWarrior,
		Tier:          1,
		Prerequisites: []string{},
		SkillPoints:   1,
		Effects: []components.SkillEffect{
			{
				Type:        "stat_bonus",
				Target:      "MaxHP",
				Value:       10,
				Description: "+10 Maximum HP",
			},
		},
		IconPath: "assets/icons/skills/tough_skin.png",
	}

	powerStrike := &components.Skill{
		ID:            "warrior_power_strike",
		Name:          "Power Strike",
		Description:   "Increases attack damage by 3 points.",
		Type:          components.SkillTypePassive,
		JobClass:      components.JobWarrior,
		Tier:          1,
		Prerequisites: []string{},
		SkillPoints:   1,
		Effects: []components.SkillEffect{
			{
				Type:        "stat_bonus",
				Target:      "Attack",
				Value:       3,
				Description: "+3 Attack Damage",
			},
		},
		IconPath: "assets/icons/skills/power_strike.png",
	}

	// Tier 2 Skills
	ironWill := &components.Skill{
		ID:            "warrior_iron_will",
		Name:          "Iron Will",
		Description:   "Increases defense and HP regeneration.",
		Type:          components.SkillTypePassive,
		JobClass:      components.JobWarrior,
		Tier:          2,
		Prerequisites: []string{"warrior_tough_skin"},
		SkillPoints:   2,
		Effects: []components.SkillEffect{
			{
				Type:        "stat_bonus",
				Target:      "Defense",
				Value:       2,
				Description: "+2 Defense",
			},
			{
				Type:        "passive_effect",
				Target:      "HP_Regen",
				Value:       1,
				Description: "1 HP regenerated per turn",
			},
		},
		IconPath: "assets/icons/skills/iron_will.png",
	}

	whirlwind := &components.Skill{
		ID:            "warrior_whirlwind",
		Name:          "Whirlwind Attack",
		Description:   "Active ability: Attack all adjacent enemies for 1 AP.",
		Type:          components.SkillTypeActive,
		JobClass:      components.JobWarrior,
		Tier:          2,
		Prerequisites: []string{"warrior_power_strike"},
		SkillPoints:   2,
		Effects: []components.SkillEffect{
			{
				Type:        "ability_unlock",
				Target:      "whirlwind_attack",
				Value:       1,
				Description: "Unlocks Whirlwind Attack ability",
				Data:        map[string]interface{}{"ap_cost": 1, "damage_multiplier": 0.8},
			},
		},
		IconPath: "assets/icons/skills/whirlwind.png",
	}

	// Register warrior skills
	sr.RegisterSkill(toughSkin)
	sr.RegisterSkill(powerStrike)
	sr.RegisterSkill(ironWill)
	sr.RegisterSkill(whirlwind)

	// Create warrior skill tree layout
	warriorTree := &components.SkillTree{
		JobClass: components.JobWarrior,
		Name:     "Warrior Combat Arts",
		Nodes:    make(map[string]*components.SkillNode),
		Layout:   make([][]*components.SkillNode, 3), // 3 rows for tiers
		MaxTier:  2,
	}

	// Create skill nodes
	toughSkinNode := &components.SkillNode{Skill: toughSkin, X: 0, Y: 0, Children: []string{"warrior_iron_will"}}
	powerStrikeNode := &components.SkillNode{Skill: powerStrike, X: 1, Y: 0, Children: []string{"warrior_whirlwind"}}
	ironWillNode := &components.SkillNode{Skill: ironWill, X: 0, Y: 1, Children: []string{}}
	whirlwindNode := &components.SkillNode{Skill: whirlwind, X: 1, Y: 1, Children: []string{}}

	// Add nodes to tree
	warriorTree.Nodes["warrior_tough_skin"] = toughSkinNode
	warriorTree.Nodes["warrior_power_strike"] = powerStrikeNode
	warriorTree.Nodes["warrior_iron_will"] = ironWillNode
	warriorTree.Nodes["warrior_whirlwind"] = whirlwindNode

	// Set up layout grid
	warriorTree.Layout = [][]*components.SkillNode{
		{toughSkinNode, powerStrikeNode}, // Tier 1
		{ironWillNode, whirlwindNode},    // Tier 2
	}

	sr.skillTrees[components.JobWarrior] = warriorTree
}

// initializeMageSkills creates the mage skill tree
func (sr *SkillRegistry) initializeMageSkills() {
	// Tier 1 Skills
	manaPool := &components.Skill{
		ID:            "mage_mana_pool",
		Name:          "Expanded Mana Pool",
		Description:   "Increases maximum MP by 15 points.",
		Type:          components.SkillTypePassive,
		JobClass:      components.JobMage,
		Tier:          1,
		Prerequisites: []string{},
		SkillPoints:   1,
		Effects: []components.SkillEffect{
			{
				Type:        "stat_bonus",
				Target:      "MaxMP",
				Value:       15,
				Description: "+15 Maximum MP",
			},
		},
		IconPath: "assets/icons/skills/mana_pool.png",
	}

	spellPower := &components.Skill{
		ID:            "mage_spell_power",
		Name:          "Spell Power",
		Description:   "Increases magical attack damage by 4 points.",
		Type:          components.SkillTypePassive,
		JobClass:      components.JobMage,
		Tier:          1,
		Prerequisites: []string{},
		SkillPoints:   1,
		Effects: []components.SkillEffect{
			{
				Type:        "stat_bonus",
				Target:      "MagicAttack",
				Value:       4,
				Description: "+4 Magic Attack",
			},
		},
		IconPath: "assets/icons/skills/spell_power.png",
	}

	// Tier 2 Skills
	fireball := &components.Skill{
		ID:            "mage_fireball",
		Name:          "Fireball",
		Description:   "Active ability: Ranged fire attack for 2 AP.",
		Type:          components.SkillTypeActive,
		JobClass:      components.JobMage,
		Tier:          2,
		Prerequisites: []string{"mage_spell_power"},
		SkillPoints:   2,
		Effects: []components.SkillEffect{
			{
				Type:        "ability_unlock",
				Target:      "fireball",
				Value:       1,
				Description: "Unlocks Fireball spell",
				Data:        map[string]interface{}{"ap_cost": 2, "range": 3, "damage": 25},
			},
		},
		IconPath: "assets/icons/skills/fireball.png",
	}

	// Register mage skills
	sr.RegisterSkill(manaPool)
	sr.RegisterSkill(spellPower)
	sr.RegisterSkill(fireball)

	// Create mage skill tree
	mageTree := &components.SkillTree{
		JobClass: components.JobMage,
		Name:     "Arcane Arts",
		Nodes:    make(map[string]*components.SkillNode),
		MaxTier:  2,
	}

	manaPoolNode := &components.SkillNode{Skill: manaPool, X: 0, Y: 0, Children: []string{}}
	spellPowerNode := &components.SkillNode{Skill: spellPower, X: 1, Y: 0, Children: []string{"mage_fireball"}}
	fireballNode := &components.SkillNode{Skill: fireball, X: 1, Y: 1, Children: []string{}}

	mageTree.Nodes["mage_mana_pool"] = manaPoolNode
	mageTree.Nodes["mage_spell_power"] = spellPowerNode
	mageTree.Nodes["mage_fireball"] = fireballNode

	mageTree.Layout = [][]*components.SkillNode{
		{manaPoolNode, spellPowerNode}, // Tier 1
		{nil, fireballNode},            // Tier 2
	}

	sr.skillTrees[components.JobMage] = mageTree
}

// initializeRogueSkills creates the rogue skill tree
func (sr *SkillRegistry) initializeRogueSkills() {
	// Tier 1 Skills
	sneak := &components.Skill{
		ID:            "rogue_sneak",
		Name:          "Sneak",
		Description:   "Move silently to avoid detection. Increases movement speed.",
		Type:          components.SkillTypePassive,
		JobClass:      components.JobRogue,
		Tier:          1,
		Prerequisites: []string{},
		SkillPoints:   1,
		Effects: []components.SkillEffect{
			{
				Type:        "stat_bonus",
				Target:      "Speed",
				Value:       3,
				Description: "+3 Speed",
			},
		},
		IconPath: "assets/icons/skills/sneak.png",
	}

	quickReflexes := &components.Skill{
		ID:            "rogue_quick_reflexes",
		Name:          "Quick Reflexes",
		Description:   "Enhanced reflexes improve defense and agility.",
		Type:          components.SkillTypePassive,
		JobClass:      components.JobRogue,
		Tier:          1,
		Prerequisites: []string{},
		SkillPoints:   1,
		Effects: []components.SkillEffect{
			{
				Type:        "stat_bonus",
				Target:      "Defense",
				Value:       2,
				Description: "+2 Defense",
			},
			{
				Type:        "stat_bonus",
				Target:      "Speed",
				Value:       2,
				Description: "+2 Speed",
			},
		},
		IconPath: "assets/icons/skills/quick_reflexes.png",
	}

	preciseStrike := &components.Skill{
		ID:            "rogue_precise_strike",
		Name:          "Precise Strike",
		Description:   "Target weak points for increased attack damage.",
		Type:          components.SkillTypePassive,
		JobClass:      components.JobRogue,
		Tier:          1,
		Prerequisites: []string{},
		SkillPoints:   1,
		Effects: []components.SkillEffect{
			{
				Type:        "stat_bonus",
				Target:      "Attack",
				Value:       4,
				Description: "+4 Attack Damage",
			},
		},
		IconPath: "assets/icons/skills/precise_strike.png",
	}

	// Tier 2 Skills
	shadowStep := &components.Skill{
		ID:            "rogue_shadow_step",
		Name:          "Shadow Step",
		Description:   "Advanced stealth techniques grant massive speed boost.",
		Type:          components.SkillTypePassive,
		JobClass:      components.JobRogue,
		Tier:          2,
		Prerequisites: []string{"rogue_sneak"},
		SkillPoints:   2,
		Effects: []components.SkillEffect{
			{
				Type:        "stat_bonus",
				Target:      "Speed",
				Value:       5,
				Description: "+5 Speed",
			},
			{
				Type:        "stat_bonus",
				Target:      "Attack",
				Value:       3,
				Description: "+3 Attack Damage",
			},
		},
		IconPath: "assets/icons/skills/shadow_step.png",
	}

	evasion := &components.Skill{
		ID:            "rogue_evasion",
		Name:          "Evasion",
		Description:   "Master dodging techniques dramatically increase survival.",
		Type:          components.SkillTypePassive,
		JobClass:      components.JobRogue,
		Tier:          2,
		Prerequisites: []string{"rogue_quick_reflexes"},
		SkillPoints:   2,
		Effects: []components.SkillEffect{
			{
				Type:        "stat_bonus",
				Target:      "Defense",
				Value:       6,
				Description: "+6 Defense",
			},
			{
				Type:        "stat_bonus",
				Target:      "MaxHP",
				Value:       8,
				Description: "+8 Maximum HP",
			},
		},
		IconPath: "assets/icons/skills/evasion.png",
	}

	deadlyStrike := &components.Skill{
		ID:            "rogue_deadly_strike",
		Name:          "Deadly Strike",
		Description:   "Devastating attack techniques for maximum damage.",
		Type:          components.SkillTypePassive,
		JobClass:      components.JobRogue,
		Tier:          2,
		Prerequisites: []string{"rogue_precise_strike"},
		SkillPoints:   2,
		Effects: []components.SkillEffect{
			{
				Type:        "stat_bonus",
				Target:      "Attack",
				Value:       7,
				Description: "+7 Attack Damage",
			},
		},
		IconPath: "assets/icons/skills/deadly_strike.png",
	}

	// Register all skills
	sr.RegisterSkill(sneak)
	sr.RegisterSkill(quickReflexes)
	sr.RegisterSkill(preciseStrike)
	sr.RegisterSkill(shadowStep)
	sr.RegisterSkill(evasion)
	sr.RegisterSkill(deadlyStrike)

	// Create the skill tree
	rogueTree := &components.SkillTree{
		JobClass: components.JobRogue,
		Name:     "Rogue Arts",
		Nodes:    make(map[string]*components.SkillNode),
		Layout:   make([][]*components.SkillNode, 3), // 3 rows for tiers
		MaxTier:  2,
	}

	// Create skill nodes
	sneakNode := &components.SkillNode{Skill: sneak, X: 0, Y: 0, Children: []string{"rogue_shadow_step"}}
	quickReflexesNode := &components.SkillNode{Skill: quickReflexes, X: 1, Y: 0, Children: []string{"rogue_evasion"}}
	preciseStrikeNode := &components.SkillNode{Skill: preciseStrike, X: 2, Y: 0, Children: []string{"rogue_deadly_strike"}}
	shadowStepNode := &components.SkillNode{Skill: shadowStep, X: 0, Y: 1, Children: []string{}}
	evasionNode := &components.SkillNode{Skill: evasion, X: 1, Y: 1, Children: []string{}}
	deadlyStrikeNode := &components.SkillNode{Skill: deadlyStrike, X: 2, Y: 1, Children: []string{}}

	// Add nodes to tree
	rogueTree.Nodes["rogue_sneak"] = sneakNode
	rogueTree.Nodes["rogue_quick_reflexes"] = quickReflexesNode
	rogueTree.Nodes["rogue_precise_strike"] = preciseStrikeNode
	rogueTree.Nodes["rogue_shadow_step"] = shadowStepNode
	rogueTree.Nodes["rogue_evasion"] = evasionNode
	rogueTree.Nodes["rogue_deadly_strike"] = deadlyStrikeNode

	// Set up layout grid
	rogueTree.Layout[0] = []*components.SkillNode{sneakNode, quickReflexesNode, preciseStrikeNode}
	rogueTree.Layout[1] = []*components.SkillNode{shadowStepNode, evasionNode, deadlyStrikeNode}
	rogueTree.Layout[2] = []*components.SkillNode{} // Empty tier 3 for now

	// Register the skill tree
	sr.skillTrees[components.JobRogue] = rogueTree
}

func (sr *SkillRegistry) initializeClericSkills() {
	// TODO: Implement cleric skills
}

func (sr *SkillRegistry) initializeArcherSkills() {
	// TODO: Implement archer skills
}

// Global skill registry instance
var GlobalSkillRegistry *SkillRegistry

// InitializeSkillRegistry initializes the global skill registry
func InitializeSkillRegistry() {
	GlobalSkillRegistry = NewSkillRegistry()
}

// GetGlobalSkillRegistry returns the global skill registry
func GetGlobalSkillRegistry() *SkillRegistry {
	if GlobalSkillRegistry == nil {
		InitializeSkillRegistry()
	}
	return GlobalSkillRegistry
}
