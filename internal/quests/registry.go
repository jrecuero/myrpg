package quests

import (
	"fmt"

	"github.com/jrecuero/myrpg/internal/ecs/components"
)

// QuestRegistry manages all available quests in the game
type QuestRegistry struct {
	quests map[string]*components.Quest // All quests keyed by ID
}

// NewQuestRegistry creates a new quest registry
func NewQuestRegistry() *QuestRegistry {
	registry := &QuestRegistry{
		quests: make(map[string]*components.Quest),
	}

	// Initialize with default quests
	registry.initializeDefaultQuests()
	return registry
}

// RegisterQuest adds a quest to the registry
func (qr *QuestRegistry) RegisterQuest(quest *components.Quest) error {
	if quest.ID == "" {
		return fmt.Errorf("quest ID cannot be empty")
	}

	if _, exists := qr.quests[quest.ID]; exists {
		return fmt.Errorf("quest with ID %s already exists", quest.ID)
	}

	qr.quests[quest.ID] = quest
	return nil
}

// GetQuest retrieves a quest by ID
func (qr *QuestRegistry) GetQuest(questID string) (*components.Quest, bool) {
	quest, exists := qr.quests[questID]
	return quest, exists
}

// GetAllQuests returns all quests
func (qr *QuestRegistry) GetAllQuests() []*components.Quest {
	quests := make([]*components.Quest, 0, len(qr.quests))
	for _, quest := range qr.quests {
		quests = append(quests, quest)
	}
	return quests
}

// GetQuestsByType returns quests of a specific type
func (qr *QuestRegistry) GetQuestsByType(questType components.QuestType) []*components.Quest {
	var typeQuests []*components.Quest
	for _, quest := range qr.quests {
		if quest.Type == questType {
			typeQuests = append(typeQuests, quest)
		}
	}
	return typeQuests
}

// GetQuestsByGiver returns quests from a specific NPC
func (qr *QuestRegistry) GetQuestsByGiver(npcID string) []*components.Quest {
	var giverQuests []*components.Quest
	for _, quest := range qr.quests {
		if quest.GiverNPCID == npcID {
			giverQuests = append(giverQuests, quest)
		}
	}
	return giverQuests
}

// CreateQuestInstance creates a new instance of a quest for a player
func (qr *QuestRegistry) CreateQuestInstance(questID string) (*components.Quest, error) {
	template, exists := qr.quests[questID]
	if !exists {
		return nil, fmt.Errorf("quest %s not found", questID)
	}

	// Create a deep copy of the quest
	questCopy := &components.Quest{
		ID:            template.ID,
		Title:         template.Title,
		Description:   template.Description,
		Type:          template.Type,
		State:         components.QuestStateInactive,
		Objectives:    make([]*components.QuestObjective, len(template.Objectives)),
		Prerequisites: append([]string{}, template.Prerequisites...),
		GiverNPCID:    template.GiverNPCID,
		Reward:        template.Reward,
		Level:         template.Level,
		Repeatable:    template.Repeatable,
	}

	// Deep copy objectives
	for i, obj := range template.Objectives {
		questCopy.Objectives[i] = &components.QuestObjective{
			ID:          obj.ID,
			Description: obj.Description,
			Type:        obj.Type,
			Target:      obj.Target,
			Required:    obj.Required,
			Current:     0, // Reset progress
			Completed:   false,
			Optional:    obj.Optional,
		}
	}

	return questCopy, nil
}

// initializeDefaultQuests creates the default quest set
func (qr *QuestRegistry) initializeDefaultQuests() {
	// Tutorial Quest - First Steps
	firstSteps := &components.Quest{
		ID:          "tutorial_first_steps",
		Title:       "First Steps",
		Description: "Learn the basics of combat and movement in this dangerous world.",
		Type:        components.QuestTypeTutorial,
		State:       components.QuestStateInactive,
		GiverNPCID:  "trainer_npc",
		Level:       1,
		Repeatable:  false,
		Objectives: []*components.QuestObjective{
			{
				ID:          "tutorial_move",
				Description: "Move around using arrow keys",
				Type:        components.ObjectiveCustom,
				Required:    1,
				Current:     0,
				Completed:   false,
				Optional:    false,
			},
			{
				ID:          "tutorial_combat",
				Description: "Defeat 2 enemies in combat",
				Type:        components.ObjectiveKill,
				Target:      "any",
				Required:    2,
				Current:     0,
				Completed:   false,
				Optional:    false,
			},
		},
		Reward: &components.QuestReward{
			Experience:  50,
			Gold:        25,
			SkillPoints: 1,
			Items:       []string{"health_potion"},
		},
	}

	// Main Quest - The Threat Emerges
	threatEmerges := &components.Quest{
		ID:            "main_threat_emerges",
		Title:         "The Threat Emerges",
		Description:   "Strange creatures have been sighted near the village. Investigate and eliminate the danger.",
		Type:          components.QuestTypeMain,
		State:         components.QuestStateInactive,
		GiverNPCID:    "village_elder",
		Level:         2,
		Repeatable:    false,
		Prerequisites: []string{"tutorial_first_steps"},
		Objectives: []*components.QuestObjective{
			{
				ID:          "eliminate_threats",
				Description: "Defeat 5 hostile creatures",
				Type:        components.ObjectiveKill,
				Target:      "enemy",
				Required:    5,
				Current:     0,
				Completed:   false,
				Optional:    false,
			},
			{
				ID:          "report_back",
				Description: "Report back to the Village Elder",
				Type:        components.ObjectiveTalk,
				Target:      "village_elder",
				Required:    1,
				Current:     0,
				Completed:   false,
				Optional:    false,
			},
		},
		Reward: &components.QuestReward{
			Experience:  100,
			Gold:        50,
			SkillPoints: 2,
			Equipment:   []string{"iron_sword"},
		},
	}

	// Side Quest - Gather Resources
	gatherResources := &components.Quest{
		ID:          "side_gather_resources",
		Title:       "Gather Resources",
		Description: "The local blacksmith needs materials for crafting. Help by collecting resources.",
		Type:        components.QuestTypeSide,
		State:       components.QuestStateInactive,
		GiverNPCID:  "blacksmith_npc",
		Level:       1,
		Repeatable:  true,
		Objectives: []*components.QuestObjective{
			{
				ID:          "collect_iron",
				Description: "Collect Iron Ore",
				Type:        components.ObjectiveCollect,
				Target:      "iron_ore",
				Required:    3,
				Current:     0,
				Completed:   false,
				Optional:    false,
			},
			{
				ID:          "collect_wood",
				Description: "Collect Wood",
				Type:        components.ObjectiveCollect,
				Target:      "wood",
				Required:    5,
				Current:     0,
				Completed:   false,
				Optional:    true, // Optional objective
			},
		},
		Reward: &components.QuestReward{
			Experience: 25,
			Gold:       30,
			Items:      []string{"crafting_materials"},
		},
	}

	// Side Quest - Equipment Mastery
	equipmentMastery := &components.Quest{
		ID:          "side_equipment_mastery",
		Title:       "Equipment Mastery",
		Description: "Learn to properly equip yourself for battle.",
		Type:        components.QuestTypeTutorial,
		State:       components.QuestStateInactive,
		GiverNPCID:  "armorer_npc",
		Level:       1,
		Repeatable:  false,
		Objectives: []*components.QuestObjective{
			{
				ID:          "equip_weapon",
				Description: "Equip a weapon",
				Type:        components.ObjectiveEquip,
				Target:      "weapon",
				Required:    1,
				Current:     0,
				Completed:   false,
				Optional:    false,
			},
			{
				ID:          "equip_armor",
				Description: "Equip armor or shield",
				Type:        components.ObjectiveEquip,
				Target:      "armor",
				Required:    1,
				Current:     0,
				Completed:   false,
				Optional:    false,
			},
		},
		Reward: &components.QuestReward{
			Experience: 30,
			Gold:       20,
			Equipment:  []string{"leather_armor"},
		},
	}

	// Daily Quest - Training Grounds
	trainingGrounds := &components.Quest{
		ID:          "daily_training",
		Title:       "Training Grounds",
		Description: "Hone your combat skills through daily practice.",
		Type:        components.QuestTypeDaily,
		State:       components.QuestStateInactive,
		GiverNPCID:  "training_master",
		Level:       1,
		Repeatable:  true,
		Objectives: []*components.QuestObjective{
			{
				ID:          "training_combat",
				Description: "Win 3 training battles",
				Type:        components.ObjectiveKill,
				Target:      "training_dummy",
				Required:    3,
				Current:     0,
				Completed:   false,
				Optional:    false,
			},
		},
		Reward: &components.QuestReward{
			Experience:  20,
			SkillPoints: 1,
		},
	}

	// Register all quests
	qr.RegisterQuest(firstSteps)
	qr.RegisterQuest(threatEmerges)
	qr.RegisterQuest(gatherResources)
	qr.RegisterQuest(equipmentMastery)
	qr.RegisterQuest(trainingGrounds)
}

// Global quest registry instance
var GlobalQuestRegistry *QuestRegistry

// InitializeQuestRegistry initializes the global quest registry
func InitializeQuestRegistry() {
	GlobalQuestRegistry = NewQuestRegistry()
}

// GetGlobalQuestRegistry returns the global quest registry
func GetGlobalQuestRegistry() *QuestRegistry {
	if GlobalQuestRegistry == nil {
		InitializeQuestRegistry()
	}
	return GlobalQuestRegistry
}
