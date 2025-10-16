# Dialog System - Character Management Guide

## Overview
This guide explains how to add, remove, and modify characters in the dialog system. Characters are defined in `assets/dialogs/characters.json` and referenced in dialog scripts.

## Character Structure

### Basic Character Definition
Each character requires the following properties:

```json
{
  "character_id": {
    "name": "Display Name",
    "portrait": "portrait_filename.png",
    "voice_style": "style_identifier",
    "description": "Character description"
  }
}
```

### Property Details
- **`character_id`**: Unique identifier used in dialog scripts (lowercase, underscores)
- **`name`**: Display name shown in dialog UI
- **`portrait`**: Portrait image filename (stored in `assets/portraits/`)
- **`voice_style`**: Style identifier for potential voice effects (currently unused)
- **`description`**: Internal description for development reference

## Adding New Characters

### Step 1: Create Character Entry
Add new character to `assets/dialogs/characters.json`:

```json
{
  "characters": {
    "existing_character": { ... },
    "new_character": {
      "name": "New Character Name",
      "portrait": "new_character_portrait.png",
      "voice_style": "appropriate_style",
      "description": "Brief character description"
    }
  }
}
```

### Step 2: Add Portrait Asset (Optional)
1. Create portrait image (recommended 64x64 pixels)
2. Save as PNG in `assets/portraits/` directory
3. Use filename specified in character definition

**Note**: If portrait file doesn't exist, dialog system shows placeholder rectangle.

### Step 3: Use Character in Dialog Scripts
Reference the character by `character_id` in dialog nodes:

```json
{
  "nodes": {
    "example_node": {
      "speaker": "new_character",
      "text": "Hello, I'm the new character!",
      "choices": [...]
    }
  }
}
```

### Example: Adding a Blacksmith
```json
"blacksmith": {
  "name": "Master Ironbeard",
  "portrait": "blacksmith_portrait.png", 
  "voice_style": "gruff",
  "description": "The village blacksmith and weapon expert"
}
```

## Removing Characters

### Step 1: Remove Character Definition
Delete the character entry from `assets/dialogs/characters.json`:

```json
{
  "characters": {
    "keep_this_character": { ... },
    // Remove this block entirely:
    // "character_to_remove": { ... }
  }
}
```

### Step 2: Update Dialog Scripts
Search all dialog files for references to the removed character:

1. **Find references**: Look for `"speaker": "character_id"`
2. **Update or remove**: Either change speaker or remove dialog nodes
3. **Test thoroughly**: Ensure no broken references remain

### Step 3: Clean Up Assets (Optional)
- Remove portrait file from `assets/portraits/` if no longer needed
- Keep organized asset directory

### ⚠️ Important Notes for Removal
- **Breaking Changes**: Removing characters breaks dialog scripts that reference them
- **Test All Dialogs**: Verify no dialog scripts reference the removed character
- **Cascade Updates**: May need to restructure conversation flows

## Modifying Existing Characters

### Changing Character Properties
Edit any property in `characters.json`:

```json
"town_elder": {
  "name": "Elder Aldric the Wise",      // ← Changed name
  "portrait": "elder_new_portrait.png", // ← New portrait
  "voice_style": "ancient",             // ← Updated style
  "description": "Ancient keeper of village wisdom" // ← New description
}
```

### Portrait Updates
1. Add new portrait file to `assets/portraits/`
2. Update `portrait` field in character definition
3. Remove old portrait file if unused

### Name Changes
- Character display name changes immediately in all dialogs
- No dialog script updates needed
- Consider narrative consistency

## Character ID Best Practices

### Naming Conventions
- **Lowercase**: `town_elder`, `merchant`, `mysterious_figure`
- **Underscores**: Use underscores for multi-word names
- **Descriptive**: Clear, memorable identifiers
- **Consistent**: Follow established patterns

### Examples
- ✅ Good: `village_guard`, `shop_keeper`, `wise_sage`
- ❌ Avoid: `Guard1`, `character_A`, `VillageGuard`

## Portrait Guidelines

### Technical Requirements
- **Format**: PNG recommended (supports transparency)
- **Size**: 64x64 pixels ideal (system handles scaling)
- **Style**: Consistent art style across characters
- **Naming**: Descriptive filenames matching character themes

### Portrait Organization
```
assets/portraits/
├── elder_portrait.png
├── merchant_portrait.png
├── guard_portrait.png
├── blacksmith_portrait.png
└── mystery_portrait.png
```

## Current Characters

### Existing Character List
| Character ID | Name | Role |
|--------------|------|------|
| `town_elder` | Elder Aldric | Village leader, quest giver |
| `merchant` | Trader Gareth | Traveling trader |
| `player` | Hero | Game protagonist |
| `guard` | Captain Marcus | Town guard captain |
| `mysterious_figure` | ??? | Enigmatic stranger |

## Integration Examples

### Adding a New Village Character
```json
// 1. Add to characters.json
"baker": {
  "name": "Baker Martha",
  "portrait": "baker_portrait.png",
  "voice_style": "warm",
  "description": "Village baker and gossip source"
}

// 2. Create dialog file: baker.json
{
  "dialog_id": "baker_chat",
  "title": "Village Baker",
  "nodes": {
    "start": {
      "speaker": "baker",
      "text": "Fresh bread today! What brings you to my bakery?",
      "choices": [...]
    }
  }
}

// 3. Use in game
g.uiManager.ShowDialog("assets/dialogs", "characters.json", "baker.json", "start")
```

### Character-Based Dialog Variations
```json
{
  "friendly_greeting": {
    "speaker": "merchant",
    "text": "Welcome, friend! Care to see my wares?"
  },
  "stern_warning": {
    "speaker": "guard", 
    "text": "Halt! State your business in this town."
  }
}
```

### Easy Dialog Endings
Use `"end": true` for automatic "End dialog" choices:

```json
{
  "brief_chat": {
    "speaker": "baker",
    "text": "Fresh bread today! Come back anytime.",
    "end": true
  }
}
```

**Result**: Shows "End dialog" button automatically - no need to manually press ESC!

## Troubleshooting

### Common Issues
1. **Missing Portrait**: Shows placeholder rectangle
   - **Solution**: Add portrait file or update filename

2. **Dialog Loading Error**: JSON syntax error
   - **Solution**: Validate JSON syntax, check trailing commas

3. **Character Not Found**: Referenced character doesn't exist
   - **Solution**: Verify character ID spelling in both files

4. **Portrait Not Loading**: File path issues
   - **Solution**: Check file exists and path is correct

### Validation Checklist
- ✅ Character ID uses lowercase and underscores
- ✅ All required properties present (name, portrait, voice_style, description)
- ✅ Portrait file exists in `assets/portraits/`
- ✅ JSON syntax is valid
- ✅ Character referenced correctly in dialog scripts
- ✅ No duplicate character IDs

## Future Enhancements

### Planned Features
- **Voice Integration**: Use `voice_style` for audio effects
- **Character Emotions**: Multiple portraits per character
- **Animation Support**: Animated character portraits
- **Character Stats**: Integration with game character data

This system provides a flexible foundation for character management that can grow with your game's narrative needs!