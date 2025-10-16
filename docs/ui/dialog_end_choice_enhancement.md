# Dialog System Enhancement - Automatic End Dialog Choice

## Overview
Enhanced the dialog system to automatically add an "End dialog" choice when a dialog node reaches an endpoint, eliminating the need to manually press ESC to close conversations.

## Implementation

### Automatic End Choice Logic
When a dialog node has:
- `"end": true` property
- No existing `choices` array (empty or missing)

The system automatically adds an "End dialog" choice that closes the dialog when selected.

### Code Changes

#### 1. Enhanced Node Processing (`displayNode` function)
```go
// Before
} else if node.End {
    // End dialog after text
    dw.State = DialogStateTyping
}

// After  
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

#### 2. Special Target Handling (`selectChoice` function)
```go
// Check for special end dialog target
if choice.Target == "__END_DIALOG__" {
    logger.Info("Dialog ended by user choice")
    dw.Hide()
    return
}
```

### Special Target System
- Uses `"__END_DIALOG__"` as internal target for automatic end choices
- Prevents conflicts with actual dialog node names
- Provides clear logging for debugging

## User Experience

### Before Enhancement
```
Elder: "May fortune favor your journey."
[Text displays]
[User must press ESC to close dialog]
```

### After Enhancement  
```
Elder: "May fortune favor your journey."
[Text displays]
┌─────────────────┐
│   End dialog    │  ← Clickable button
└─────────────────┘
```

## Behavior Rules

### When Automatic Choice Appears
✅ Node has `"end": true`  
✅ Node has no `choices` array or empty array  
✅ Appears after typewriter text completes

### When Automatic Choice Does NOT Appear  
❌ Node has existing `choices` array with content  
❌ Node has `"next"` auto-advance property  
❌ Node does not have `"end": true`

### Priority System
1. **Existing choices** - If node has choices, use them
2. **Automatic end choice** - If node has `"end": true` and no choices  
3. **Auto-advance** - If node has `"next"` property
4. **Wait for input** - Default behavior

## Documentation Updates

### Updated Files
- `docs/ui/dialog_script_format.md` - Added "Dialog Ending Behavior" section
- `docs/ui/dialog_character_management.md` - Added example in Integration section

### Key Documentation Points
- Explains automatic vs manual end choices
- Shows examples of both approaches
- Clarifies priority rules for mixed scenarios

## Examples

### Simple End Dialog
```json
{
  "farewell": {
    "speaker": "town_elder",
    "text": "Safe travels, adventurer.",
    "end": true
  }
}
```
**Result**: Shows "End dialog" button after text

### Manual End Choices (Still Supported)
```json
{
  "conversation": {
    "speaker": "merchant", 
    "text": "Anything else I can help with?",
    "choices": [
      {
        "text": "Show me your wares.",
        "target": "shop_menu"
      },
      {
        "text": "Goodbye.",
        "target": "farewell" 
      }
    ]
  }
}
```
**Result**: Shows custom choices, no automatic end choice

## Benefits

### User Experience
- ✅ **Intuitive**: Clear "End dialog" button instead of hidden ESC key
- ✅ **Consistent**: Same interaction pattern as other choices
- ✅ **Accessible**: Mouse-clickable like other dialog options
- ✅ **Visual**: Clear endpoint indication

### Developer Experience  
- ✅ **Simple**: Just add `"end": true` to dialog nodes
- ✅ **Flexible**: Still supports manual choice control when needed
- ✅ **Backward Compatible**: Existing dialogs work unchanged
- ✅ **Maintainable**: Clean internal implementation with special targets

## Testing

### Verification Steps
1. ✅ Build compiles successfully
2. ✅ Dialog test runs without errors  
3. ✅ JSON validation passes for all dialog files
4. ✅ Automatic end choices appear correctly
5. ✅ End choice properly closes dialog
6. ✅ Manual choices still work as before

### Test Cases
- End nodes with `"end": true` and no choices → Shows "End dialog"
- End nodes with existing choices → Uses existing choices  
- Non-end nodes → No automatic end choice
- Mixed dialog trees → Proper choice handling throughout

## Future Enhancements

### Potential Improvements
- **Customizable Text**: Allow custom end choice text per dialog
- **Multiple End Types**: Different end actions (save, return to menu, etc.)
- **Confirmation Dialogs**: "Are you sure you want to end?" for important conversations
- **Keyboard Shortcuts**: Still support ESC while showing the button

This enhancement significantly improves the dialog system's usability while maintaining full backward compatibility and developer flexibility!