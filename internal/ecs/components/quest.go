package components

import (
	"fmt"
	"time"
)

// QuestState represents the current state of a quest
type QuestState int

const (
	QuestStateInactive  QuestState = iota // Quest not yet available/discovered
	QuestStateActive                      // Quest is currently active
	QuestStateCompleted                   // Quest has been completed
	QuestStateFailed                      // Quest has failed (optional)
)

// QuestType categorizes different types of quests
type QuestType int

const (
	QuestTypeMain     QuestType = iota // Main story quest
	QuestTypeSide                      // Side quest
	QuestTypeDaily                     // Daily/repeatable quest
	QuestTypeTutorial                  // Tutorial quest
)

// ObjectiveType defines different types of quest objectives
type ObjectiveType int

const (
	ObjectiveKill    ObjectiveType = iota // Kill X enemies
	ObjectiveCollect                      // Collect X items
	ObjectiveTalk                         // Talk to NPC
	ObjectiveVisit                        // Visit location
	ObjectiveEquip                        // Equip item
	ObjectiveLevel                        // Reach level X
	ObjectiveCustom                       // Custom objective with manual completion
)

// QuestObjective represents a single objective within a quest
type QuestObjective struct {
	ID          string        `json:"id"`
	Description string        `json:"description"`
	Type        ObjectiveType `json:"type"`
	Target      string        `json:"target"`    // What to interact with (enemy type, item ID, NPC ID, etc.)
	Required    int           `json:"required"`  // How many needed
	Current     int           `json:"current"`   // Current progress
	Completed   bool          `json:"completed"` // Whether this objective is complete
	Optional    bool          `json:"optional"`  // Whether this objective is optional
}

// IsCompleted checks if the objective is finished
func (obj *QuestObjective) IsCompleted() bool {
	return obj.Completed || (obj.Current >= obj.Required)
}

// UpdateProgress increments the objective progress
func (obj *QuestObjective) UpdateProgress(amount int) {
	obj.Current += amount
	if obj.Current >= obj.Required {
		obj.Current = obj.Required
		obj.Completed = true
	}
}

// GetProgressText returns a formatted progress string
func (obj *QuestObjective) GetProgressText() string {
	if obj.Type == ObjectiveTalk || obj.Type == ObjectiveVisit || obj.Type == ObjectiveCustom {
		if obj.Completed {
			return "✓ " + obj.Description
		}
		return "○ " + obj.Description
	}

	status := "○"
	if obj.Completed {
		status = "✓"
	}
	return fmt.Sprintf("%s %s (%d/%d)", status, obj.Description, obj.Current, obj.Required)
}

// QuestReward represents rewards given upon quest completion
type QuestReward struct {
	Experience  int      `json:"experience"`   // XP reward
	Gold        int      `json:"gold"`         // Gold reward
	Items       []string `json:"items"`        // Item IDs to give
	SkillPoints int      `json:"skill_points"` // Skill points to award
	Equipment   []string `json:"equipment"`    // Equipment IDs to give
}

// Quest represents a complete quest definition
type Quest struct {
	ID            string            `json:"id"`
	Title         string            `json:"title"`
	Description   string            `json:"description"`
	Type          QuestType         `json:"type"`
	State         QuestState        `json:"state"`
	Objectives    []*QuestObjective `json:"objectives"`
	Prerequisites []string          `json:"prerequisites"` // Quest IDs that must be completed first
	GiverNPCID    string            `json:"giver_npc_id"`  // NPC who gives this quest
	Reward        *QuestReward      `json:"reward"`
	StartedAt     time.Time         `json:"started_at"`
	CompletedAt   time.Time         `json:"completed_at"`
	Level         int               `json:"level"`      // Recommended level
	Repeatable    bool              `json:"repeatable"` // Can be repeated
}

// CanStart checks if the quest can be started
func (q *Quest) CanStart() bool {
	return q.State == QuestStateInactive
}

// IsActive checks if the quest is currently active
func (q *Quest) IsActive() bool {
	return q.State == QuestStateActive
}

// IsCompleted checks if the quest is completed
func (q *Quest) IsCompleted() bool {
	return q.State == QuestStateCompleted
}

// GetCompletedObjectives returns the number of completed objectives
func (q *Quest) GetCompletedObjectives() int {
	completed := 0
	for _, obj := range q.Objectives {
		if obj.IsCompleted() {
			completed++
		}
	}
	return completed
}

// GetRequiredObjectives returns the number of required (non-optional) objectives
func (q *Quest) GetRequiredObjectives() int {
	required := 0
	for _, obj := range q.Objectives {
		if !obj.Optional {
			required++
		}
	}
	return required
}

// CheckCompletion checks if all required objectives are complete
func (q *Quest) CheckCompletion() bool {
	for _, obj := range q.Objectives {
		if !obj.Optional && !obj.IsCompleted() {
			return false
		}
	}
	return true
}

// StartQuest activates the quest
func (q *Quest) StartQuest() {
	q.State = QuestStateActive
	q.StartedAt = time.Now()
}

// CompleteQuest marks the quest as completed
func (q *Quest) CompleteQuest() {
	q.State = QuestStateCompleted
	q.CompletedAt = time.Now()

	// Mark all objectives as completed
	for _, obj := range q.Objectives {
		if !obj.Optional {
			obj.Completed = true
		}
	}
}

// UpdateObjective updates progress on a specific objective
func (q *Quest) UpdateObjective(objectiveID string, amount int) bool {
	for _, obj := range q.Objectives {
		if obj.ID == objectiveID {
			wasCompleted := obj.IsCompleted()
			obj.UpdateProgress(amount)
			return !wasCompleted && obj.IsCompleted() // Return true if objective was just completed
		}
	}
	return false
}

// QuestJournalComponent manages a character's quest progress
type QuestJournalComponent struct {
	ActiveQuests    map[string]*Quest `json:"active_quests"`    // Currently active quests
	CompletedQuests map[string]*Quest `json:"completed_quests"` // Completed quests
	FailedQuests    map[string]*Quest `json:"failed_quests"`    // Failed quests (optional)
}

// NewQuestJournalComponent creates a new quest journal
func NewQuestJournalComponent() *QuestJournalComponent {
	return &QuestJournalComponent{
		ActiveQuests:    make(map[string]*Quest),
		CompletedQuests: make(map[string]*Quest),
		FailedQuests:    make(map[string]*Quest),
	}
}

// AddQuest adds a new quest to the journal
func (qj *QuestJournalComponent) AddQuest(quest *Quest) {
	if quest.IsActive() {
		qj.ActiveQuests[quest.ID] = quest
	} else if quest.IsCompleted() {
		qj.CompletedQuests[quest.ID] = quest
	}
}

// StartQuest activates a quest
func (qj *QuestJournalComponent) StartQuest(questID string, quest *Quest) {
	quest.StartQuest()
	qj.ActiveQuests[questID] = quest
}

// CompleteQuest marks a quest as completed and moves it
func (qj *QuestJournalComponent) CompleteQuest(questID string) *QuestReward {
	if quest, exists := qj.ActiveQuests[questID]; exists {
		quest.CompleteQuest()
		qj.CompletedQuests[questID] = quest
		delete(qj.ActiveQuests, questID)
		return quest.Reward
	}
	return nil
}

// GetActiveQuests returns all active quests
func (qj *QuestJournalComponent) GetActiveQuests() []*Quest {
	quests := make([]*Quest, 0, len(qj.ActiveQuests))
	for _, quest := range qj.ActiveQuests {
		quests = append(quests, quest)
	}
	return quests
}

// GetCompletedQuests returns all completed quests
func (qj *QuestJournalComponent) GetCompletedQuests() []*Quest {
	quests := make([]*Quest, 0, len(qj.CompletedQuests))
	for _, quest := range qj.CompletedQuests {
		quests = append(quests, quest)
	}
	return quests
}

// UpdateQuestProgress updates progress on quest objectives
func (qj *QuestJournalComponent) UpdateQuestProgress(objectiveType ObjectiveType, target string, amount int) []string {
	var completedObjectives []string

	for _, quest := range qj.ActiveQuests {
		for _, obj := range quest.Objectives {
			if obj.Type == objectiveType && obj.Target == target && !obj.IsCompleted() {
				wasCompleted := obj.IsCompleted()
				obj.UpdateProgress(amount)
				if !wasCompleted && obj.IsCompleted() {
					completedObjectives = append(completedObjectives, quest.Title+" - "+obj.Description)
				}
			}
		}
	}

	return completedObjectives
}

// HasQuest checks if a quest is in the journal (active or completed)
func (qj *QuestJournalComponent) HasQuest(questID string) bool {
	_, activeExists := qj.ActiveQuests[questID]
	_, completedExists := qj.CompletedQuests[questID]
	return activeExists || completedExists
}

// GetQuest retrieves a quest by ID from active or completed
func (qj *QuestJournalComponent) GetQuest(questID string) *Quest {
	if quest, exists := qj.ActiveQuests[questID]; exists {
		return quest
	}
	if quest, exists := qj.CompletedQuests[questID]; exists {
		return quest
	}
	return nil
}
