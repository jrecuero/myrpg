// Package save provides game state persistence functionality.
package save

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/events"
	"github.com/jrecuero/myrpg/internal/logger"
)

// SaveManager handles game state persistence
type SaveManager struct {
	saveDirectory string              // Directory where save files are stored
	eventSaveData *EventStateSaveData // Current event save data
}

// NewSaveManager creates a new save manager
func NewSaveManager(saveDirectory string) *SaveManager {
	// Create save directory if it doesn't exist
	if err := os.MkdirAll(saveDirectory, 0755); err != nil {
		logger.Error("Failed to create save directory %s: %v", saveDirectory, err)
		return nil
	}

	return &SaveManager{
		saveDirectory: saveDirectory,
		eventSaveData: NewEventStateSaveData(),
	}
}

// SaveEventState saves event state for all entities with event components
func (sm *SaveManager) SaveEventState(world *ecs.World, eventManager *events.EventManager) error {
	if sm.eventSaveData == nil {
		sm.eventSaveData = NewEventStateSaveData()
	}

	logger.Info("Saving event state...")

	// Collect event states from all entities
	entityCount := 0
	eventCount := 0

	for _, entity := range world.GetEntities() {
		if ec := entity.Event(); ec != nil {
			sm.eventSaveData.AddEventState(ec)
			eventCount++
		}
		entityCount++
	}

	// Also save the completion tracking from event manager
	if eventManager != nil {
		// Get current game mode from event manager if available
		// Note: This would require adding a GetGameMode method to EventManager
		// For now, we'll keep the current mode in the save data
	}

	// Generate save file path with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("events_%s.json", timestamp)
	filePath := filepath.Join(sm.saveDirectory, filename)

	// Create also a "latest" save for easy loading
	latestPath := filepath.Join(sm.saveDirectory, "events_latest.json")

	// Save to JSON
	data, err := sm.eventSaveData.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize event data: %v", err)
	}

	// Write to timestamped file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write save file %s: %v", filePath, err)
	}

	// Write to latest file
	if err := os.WriteFile(latestPath, data, 0644); err != nil {
		logger.Warn("Failed to write latest save file %s: %v", latestPath, err)
	}

	logger.Info("Event state saved successfully: %d events from %d entities to %s", eventCount, entityCount, filePath)
	return nil
}

// LoadEventState loads event state from the latest save file
func (sm *SaveManager) LoadEventState() (*EventStateSaveData, error) {
	latestPath := filepath.Join(sm.saveDirectory, "events_latest.json")

	// Check if latest save exists
	if _, err := os.Stat(latestPath); os.IsNotExist(err) {
		logger.Info("No event save file found at %s", latestPath)
		return NewEventStateSaveData(), nil // Return empty save data
	}

	// Read the file
	data, err := os.ReadFile(latestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read save file %s: %v", latestPath, err)
	}

	// Deserialize
	saveData := NewEventStateSaveData()
	if err := saveData.FromJSON(data); err != nil {
		return nil, fmt.Errorf("failed to parse save file %s: %v", latestPath, err)
	}

	// Validate the loaded data
	if err := saveData.Validate(); err != nil {
		return nil, fmt.Errorf("save file validation failed: %v", err)
	}

	logger.Info("Event state loaded successfully: %d events, %d completed from %s",
		saveData.GetTotalEventCount(), saveData.GetCompletedEventCount(), latestPath)

	sm.eventSaveData = saveData
	return saveData, nil
}

// LoadEventStateFromFile loads event state from a specific file
func (sm *SaveManager) LoadEventStateFromFile(filename string) (*EventStateSaveData, error) {
	filePath := filepath.Join(sm.saveDirectory, filename)

	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read save file %s: %v", filePath, err)
	}

	// Deserialize
	saveData := NewEventStateSaveData()
	if err := saveData.FromJSON(data); err != nil {
		return nil, fmt.Errorf("failed to parse save file %s: %v", filePath, err)
	}

	// Validate the loaded data
	if err := saveData.Validate(); err != nil {
		return nil, fmt.Errorf("save file validation failed: %v", err)
	}

	logger.Info("Event state loaded from %s: %d events, %d completed",
		filename, saveData.GetTotalEventCount(), saveData.GetCompletedEventCount())

	return saveData, nil
}

// ApplyEventState applies loaded event state to the world and event manager
func (sm *SaveManager) ApplyEventState(world *ecs.World, eventManager *events.EventManager, saveData *EventStateSaveData) error {
	if saveData == nil {
		return fmt.Errorf("save data is nil")
	}

	logger.Info("Applying event state to world...")

	appliedCount := 0
	notFoundCount := 0

	// Apply state to all entities with event components
	for _, entity := range world.GetEntities() {
		if ec := entity.Event(); ec != nil {
			if savedState, exists := saveData.GetEventState(ec.ID); exists {
				savedState.ApplyToEventComponent(ec)
				appliedCount++
				logger.Debug("Applied saved state to event %s: state=%v, triggers=%d",
					ec.ID, ec.State, ec.TriggerCount)
			} else {
				notFoundCount++
				logger.Debug("No saved state found for event %s, keeping defaults", ec.ID)
			}
		}
	}

	// Update event manager completion tracking if available
	if eventManager != nil {
		for eventID, completed := range saveData.CompletedEvents {
			if completed {
				// Note: This would require adding a SetEventCompleted method to EventManager
				// For now, we'll log the information
				logger.Debug("Event %s should be marked as completed", eventID)
			}
		}
	}

	logger.Info("Event state applied: %d events updated, %d events not found in save",
		appliedCount, notFoundCount)

	return nil
}

// ListSaveFiles returns a list of available save files
func (sm *SaveManager) ListSaveFiles() ([]string, error) {
	entries, err := os.ReadDir(sm.saveDirectory)
	if err != nil {
		return nil, fmt.Errorf("failed to read save directory: %v", err)
	}

	var saveFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			saveFiles = append(saveFiles, entry.Name())
		}
	}

	return saveFiles, nil
}

// DeleteSaveFile deletes a specific save file
func (sm *SaveManager) DeleteSaveFile(filename string) error {
	filePath := filepath.Join(sm.saveDirectory, filename)

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete save file %s: %v", filePath, err)
	}

	logger.Info("Deleted save file: %s", filename)
	return nil
}

// GetSaveFileInfo returns information about a save file
func (sm *SaveManager) GetSaveFileInfo(filename string) (map[string]interface{}, error) {
	saveData, err := sm.LoadEventStateFromFile(filename)
	if err != nil {
		return nil, err
	}

	info := map[string]interface{}{
		"filename":   filename,
		"statistics": saveData.GetStatistics(),
	}

	// Add file system info
	filePath := filepath.Join(sm.saveDirectory, filename)
	if stat, err := os.Stat(filePath); err == nil {
		info["file_size"] = stat.Size()
		info["modified"] = stat.ModTime()
	}

	return info, nil
}

// GetCurrentEventSaveData returns the current event save data
func (sm *SaveManager) GetCurrentEventSaveData() *EventStateSaveData {
	return sm.eventSaveData
}

// ClearEventState resets all event states to defaults
func (sm *SaveManager) ClearEventState(world *ecs.World) error {
	logger.Info("Clearing all event states to defaults...")

	resetCount := 0

	for _, entity := range world.GetEntities() {
		if ec := entity.Event(); ec != nil {
			ec.Reset() // Use the built-in Reset method
			resetCount++
			logger.Debug("Reset event %s to default state", ec.ID)
		}
	}

	// Clear the save data
	sm.eventSaveData = NewEventStateSaveData()

	logger.Info("Cleared %d events to default state", resetCount)
	return nil
}

// GetSaveDirectory returns the save directory path
func (sm *SaveManager) GetSaveDirectory() string {
	return sm.saveDirectory
}
