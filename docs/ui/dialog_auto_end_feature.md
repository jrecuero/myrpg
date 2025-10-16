# Dialog System - Auto End Dialog Feature

## Overview
The dialog system automatically adds an "End dialog" choice to dialog nodes that have `"end": true` but no explicit choices defined. This provides a user-friendly way to close dialogs without requiring the ESC key.

## How It Works

### Automatic End Choice Generation
When a dialog node meets these conditions:
1. Has `"end": true` property
2. Has no `choices` array or empty choices

The system automatically adds:
```
Choice: "End dialog"
```

### User Experience
- Users see "End dialog" as a clickable option
- Clicking it immediately closes the dialog widget
- No need to remember the ESC key
- Consistent with other choice-based navigation

## Implementation Details

### Code Location
File: `internal/ui/dialog_widget.go`

### Key Functions

#### Node Processing (lines 354-362)
```go
} else if node.End {
    // Add automatic "End dialog" choice for end nodes
    dw.State = DialogStateTyping
    endChoice := Choice{
        Text:   "End dialog",
        Target: "__END_DIALOG__", // Special target to close dialog
    }
    dw.AvailableChoices = []Choice{endChoice}
    dw.SelectedChoice = 0
    dw.ChoiceScrollOffset = 0
}
```

#### Choice Selection (lines 640-644)
```go
if choice.Target == "__END_DIALOG__" {
    logger.Info("Dialog ended by user choice")
    dw.Hide()
    return
}
```

### Special Target ID
- **`__END_DIALOG__`**: Reserved target that triggers dialog closure
- Not a real dialog node - handled specially in choice processing
- Cannot be used as a regular node ID

## Dialog Script Usage

### End Nodes Without Choices
```json
{
  "farewell_node": {
    "speaker": "town_elder",
    "text": "May your journey be safe, traveler.",
    "end": true
  }
}
```
**Result**: Shows "End dialog" choice automatically

### End Nodes With Explicit Choices
```json
{
  "quest_complete": {
    "speaker": "town_elder", 
    "text": "You have completed the quest!",
    "end": true,
    "choices": [
      {
        "text": "Thank you for the reward.",
        "target": "end_grateful"
      },
      {
        "text": "I must continue my journey.",
        "target": "end_departure"
      }
    ]
  }
}
```
**Result**: Shows explicit choices, no auto "End dialog" added

### Non-End Nodes
```json
{
  "middle_conversation": {
    "speaker": "merchant",
    "text": "What would you like to know?",
    "choices": [
      {"text": "About your wares", "target": "show_items"},
      {"text": "About the road ahead", "target": "road_info"}
    ]
  }
}
```
**Result**: Shows only explicit choices, no auto end option

## Current Dialog Files

### Town Elder Dialog
Contains multiple end nodes that will show "End dialog":
- `end_friendly` 
- `end_helpful`
- `end_quest_given`
- `end_preparation`

### Merchant Dialog  
Contains end nodes that will show "End dialog":
- `end_transaction`
- `end_polite`
- `end_farewell`

## Benefits

### User Experience
- **Intuitive**: Clear indication dialog can be closed
- **Consistent**: All dialogs have same ending mechanism
- **Accessible**: No need to know keyboard shortcuts
- **Visual**: Matches other choice-based interactions

### Development
- **Automatic**: No need to manually add end choices
- **Flexible**: Can still use explicit choices when desired
- **Backward Compatible**: Existing dialogs work unchanged
- **Consistent**: Standardized dialog ending behavior

## Testing

### In Dialog Test Program
1. Run `make test-dialog`
2. Press 'D' to start elder dialog
3. Navigate through conversation choices
4. Reach any ending node (e.g., "Farewell for now.")
5. Observe automatic "End dialog" choice appears
6. Click "End dialog" to close

### In Main Game
1. Press 'D' in main game
2. Navigate through dialog
3. Reach ending nodes
4. Use "End dialog" choice to close

## Technical Notes

### State Management
- End choice treated like regular choice in UI
- Uses same selection/navigation system
- Proper cleanup when dialog closes

### Logging
- Dialog closure logged: "Dialog ended by user choice"
- Helps distinguish user choice vs ESC key closure

### Edge Cases
- Nodes with `"end": true` AND choices: No auto end added (explicit choices take precedence)
- Nodes without `"end": true`: No auto end added (continue conversation)
- Empty choices array: Treated same as no choices (auto end added if `"end": true`)

This feature makes dialog navigation more intuitive while maintaining full backward compatibility with existing dialog scripts!